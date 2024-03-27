package main

import (
	_ "github.com/mattn/go-sqlite3"

	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/evmos/validator-status/pkg/config"
	"github.com/evmos/validator-status/pkg/database"
	"github.com/evmos/validator-status/pkg/logger"
	"github.com/evmos/validator-status/pkg/server"
	"github.com/evmos/validator-status/sql"
)

func main() {
	configFile := flag.String("c", "./config.json", "path to config file")
	flag.Parse()

	// Log to file
	conf := config.FromFile(*configFile)
	if conf.LogFile != "" {
		logFile := logger.LogToFile(conf.LogFile)
		defer logFile.Close()
	}
	logger.SetFlags()

	// Handle signals
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	// Init db
	db, err := sql.InitDatabase(context.Background(), conf)
	if err != nil {
		logger.LogInfo(fmt.Sprintf("failed to start the database: %s", err.Error()))
		panic(err)
	}
	queries := database.New(db)

	// Run the server
	logger.LogInfo("starting the program")
	s := server.NewServer(conf, queries)
	s.Run()
	defer s.Shutdown()

	// Stop the server
	<-stop
	logger.LogInfo("stopping the program")
}
