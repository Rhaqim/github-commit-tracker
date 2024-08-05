package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Rhaqim/savannahtech/old/core"
	"github.com/Rhaqim/savannahtech/old/database"
	"github.com/Rhaqim/savannahtech/old/model"
	"github.com/Rhaqim/savannahtech/old/router"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func main() {
	// Initialise the logger
	// log.Init(true)

	// Initialize the database
	// database.Init()
	// defer database.Close()

	// Run database migrations
	err := model.Migrations()
	if err != nil {
		logger.ErrorLogger.Fatal("Failed to run database migrations:", err)
	}

	// Initialize the cache
	database.CacheInit()
	defer database.CacheClose()

	// Create a channel to handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the event listeners
	EventListeners()

	// load the startup repo
	if err := core.LoadStartupRepo(); err != nil {
		logger.ErrorLogger.Fatal("Failed to load startup repo:", err)
	}

	// Start the Gin server
	go func() {
		r := router.NewRouter()
		if err := r.Run(":8080"); err != nil {
			logger.ErrorLogger.Fatalf("Failed to run Gin server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	<-signalChan
	logger.InfoLogger.Println("Shutting down server and event listener...")
	// Perform any cleanup here, like closing database connections or gracefully stopping event listeners
	os.Exit(0)

}

func EventListeners() {
	// Start the event listener for repo events
	go func() {
		if err := core.GetRepoEvent(); err != nil {
			logger.ErrorLogger.Fatalf("Error in event listener: %v", err)
		}
	}()

	// Start the event listener for repo events
	go func() {
		if err := core.GetCommitEvent(); err != nil {
			logger.ErrorLogger.Fatalf("Error in event listener: %v", err)
		}
	}()

	// start the periodic fetch
	go func() {
		if err := core.PeriodFetch(); err != nil {
			logger.ErrorLogger.Fatalf("Error in periodic fetch: %v", err)
		}
	}()
}
