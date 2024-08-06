package main

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/internal/database"
	"github.com/Rhaqim/savannahtech/internal/listener"
	"github.com/Rhaqim/savannahtech/internal/router"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func main() {
	// Initialize logger
	logger.InitLogger(true)

	// Load environment variables
	config.LoadConfig()

	// Initialize database
	database.InitDB()
	defer database.Close()

	// Run database migrations
	database.RunMigrations()

	go func() {
		// Set up the router
		r := router.SetupRouter()

		// Start the server
		r.Run(config.Config.ServerPort)
	}()

	// Create a channel to handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	// Initialize and start event listeners
	// listener.EventListeners()
	go listener.StartEventListeners()

	// Load default
	services.LoadStartupRepo()

	<-signalChan
	logger.InfoLogger.Println("Shutting down server and event listener...")
	// Perform any cleanup here, like closing database connections or gracefully stopping event listeners
	os.Exit(0)

}
