package events

import (
	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

var (
	EventChan = make(chan entities.Event)
	ErrorChan = make(chan error, 1)
)

func startEventListeners(processFunc func(owner, repo, start_date string) error) {
	go func() {
		for event := range EventChan {
			logger.InfoLogger.Printf("Event received: %v", event)

			if err := processFunc(event.Owner, event.Repo, event.From); err != nil {
				ErrorChan <- err
			}
		}
	}()

	// Handle errors from the ErrorChan
	go func() {
		for err := range ErrorChan {
			if err != nil {
				logger.ErrorLogger.Printf("Error processing event: %v", err)
			}
		}
	}()

	logger.InfoLogger.Println("Event listeners started")
}

func SendEvent(event entities.Event) {
	EventChan <- event
}

func StartCommitEventListener() {
	startEventListeners(services.ProcessCommitData)
}

// func StartPeriodicFetchListener() {
// 	startEventListeners(PeriodicFetch)
// }

func StartEventListeners() {
	StartCommitEventListener()
	// StartPeriodicFetchListener()
}
