package server

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/evmos/validator-status/pkg/config"
	"github.com/evmos/validator-status/pkg/cosmos"
	"github.com/evmos/validator-status/pkg/database"
	"github.com/gorilla/mux"
)

type Server struct {
	config *config.Config
	Router *mux.Router
	Server *http.Server

	// Go rutine to run updates
	updateTicker *time.Ticker
	doneTicker   chan bool

	// Database
	db *database.Queries

	// Cosmos
	cosmos *cosmos.Cosmos
}

func NewServer(config *config.Config, db *database.Queries) *Server {
	router := mux.NewRouter()
	EnableCORS(router)

	tickerDuration, err := time.ParseDuration(config.RefreshDuration)
	if err != nil {
		panic(fmt.Sprintf("invalid refresh duration %q", err))
	}

	server := &Server{
		config:       config,
		Router:       router,
		Server:       &http.Server{},
		updateTicker: time.NewTicker(tickerDuration),
		doneTicker:   make(chan bool),
		db:           db,
		cosmos:       cosmos.NewCosmos(config.CosmosAPI, db),
	}

	router.Path(HomePrefix).HandlerFunc(server.HandlerHome).Methods(http.MethodGet, http.MethodPost, http.MethodOptions)

	return server
}

func (s *Server) Run() {
	if s.config.HTTPServer.Port == "" {
		panic("PORT was not set")
	}
	if s.config.HTTPServer.Address == "" {
		s.config.HTTPServer.Address = "localhost"
	}

	s.Server = &http.Server{
		Addr:              fmt.Sprintf("%s:%s", s.config.HTTPServer.Address, s.config.HTTPServer.Port),
		Handler:           s.Router,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go s.RunUpdates()
	go func() {
		_ = s.Server.ListenAndServe()
	}()
}

func (s *Server) Shutdown() {
	s.StopUpdates()
	if err := s.Server.Shutdown(context.Background()); err != nil {
		panic(fmt.Sprintf("error shuting down the server %q", err))
	}
}
