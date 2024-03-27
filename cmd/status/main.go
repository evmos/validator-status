package main

import (
	_ "github.com/mattn/go-sqlite3"

	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/evmos/validator-status/pkg/config"
	"github.com/evmos/validator-status/pkg/database"
	"github.com/evmos/validator-status/pkg/logger"
	"github.com/evmos/validator-status/pkg/server"
	"github.com/evmos/validator-status/sql"
)

func main() {
	configFile := flag.String("c", "./config.json", "path to config file")
	flag.Parse()

	conf := config.FromFile(*configFile)

	if conf.LogFile != "" {
		logFile := logger.LogToFile(conf.LogFile)
		defer logFile.Close()
	}
	logger.SetFlags()

	run := true
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		logger.LogInfo("stopping the program")
		run = false
	}()

	db, err := sql.InitDatabase(context.Background(), conf)
	if err != nil {
		logger.LogInfo(fmt.Sprintf("failed to start the database: %s", err.Error()))
		panic(err)
	}

	queries := database.New(db)

	logger.LogInfo("starting the server")
	s := server.NewServer(conf, queries)
	s.Run()
	defer s.Shutdown()

	for run {
		time.Sleep(100 * time.Millisecond)
	}

}
