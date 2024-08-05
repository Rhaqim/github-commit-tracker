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
		r.Run(config.Config.ServerAddress)
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

/*
Clean Code
- Delete unused code
- The name of the struct variable in struct methods should be a lowercase letter
- The name of the module should contain the full GitHub path (github.com/rhaqim/savannah)
- Using string concatenation is not conventional, using fmt.Sprintf is better (postgres.go, internal/core/repo.go, periodic.go, internal/api/repo.go)
- Unit tests should be located in the same folder as the functions themselves (filename: FILENAME_test.go)
- The code is not modular, the data layer is intertwined with the code, and therefore unit tests can’t be created for any data-related function, a minimal interface for the data layer is missing
- The gorm implementation is problematic. You should be using a repository to access your persistent layer.
- migration.go should be under database folder
- errChan in startEventListener function is redundant, all it is used for is printing the errors, they can be printed in the processFunc goroutine

Code Flow
- Instead of redis pub sub can you use channels to handle the communications?
- godotenv.Load should be executed in main.go and your config should be a singleton object
- The service should wait for the event listeners to initiate successfully
- Internal server errors should be logged, but should not reach the client
- HTTP response to a GET request should be included in the response body, and not in a header

Data
- Commit_url column is redundant as it can be built by the two other columns
- Missing index for author request
*/
