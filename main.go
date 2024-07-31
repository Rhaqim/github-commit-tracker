package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"savannahtech/src/core"
	"savannahtech/src/database"
	"savannahtech/src/model"
	"savannahtech/src/router"
)

func main() {
	database.Init()
	defer database.Close()

	err := model.Migrations()
	if err != nil {
		log.Fatal(err)
	}

	// Initialize the cache
	database.CacheInit()
	defer database.CacheClose()

	// Create a channel to handle OS signals for graceful shutdown
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)

	go func() {
		r := router.NewRouter()
		if err := r.Run(":8080"); err != nil {
			log.Fatalf("Failed to run Gin server: %v", err)
		}
	}()

	// Start the event listener in a separate goroutine
	go func() {
		if err := core.GetEvent(); err != nil {
			log.Fatalf("Error in event listener: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the application
	<-signalChan
	log.Println("Shutting down server and event listener...")
	// Perform any cleanup here, like closing database connections or gracefully stopping event listeners

}
