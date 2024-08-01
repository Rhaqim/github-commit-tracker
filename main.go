package main

import (
	"os"
	"os/signal"
	"syscall"

	"savannahtech/src/core"
	"savannahtech/src/database"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/router"
)

func main() {
	// Initialise the logger
	log.Init()

	// Initialize the database
	database.Init()
	defer database.Close()

	// Run database migrations
	err := model.Migrations()
	if err != nil {
		log.ErrorLogger.Fatal("Failed to run database migrations:", err)
	}

	// Initialize the cache
	database.CacheInit()
	defer database.CacheClose()

	// Create a channel to handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Start the event listeners
	EventListeners()

	// Start the Gin server
	go func() {
		r := router.NewRouter()
		if err := r.Run(":8080"); err != nil {
			log.ErrorLogger.Fatalf("Failed to run Gin server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	<-signalChan
	log.InfoLogger.Println("Shutting down server and event listener...")
	// Perform any cleanup here, like closing database connections or gracefully stopping event listeners
	os.Exit(0)

}

func EventListeners() {
	// Start the event listener for repo events
	go func() {
		if err := core.GetCommitEvent(); err != nil {
			log.ErrorLogger.Fatalf("Error in event listener: %v", err)
		}
	}()

	// Start the event listener for repo events
	go func() {
		if err := core.GetRepoEvent(); err != nil {
			log.ErrorLogger.Fatalf("Error in event listener: %v", err)
		}
	}()

	// start the periodic fetch
	go func() {
		if err := core.PeriodFetch(); err != nil {
			log.ErrorLogger.Fatalf("Error in periodic fetch: %v", err)
		}
	}()
}
