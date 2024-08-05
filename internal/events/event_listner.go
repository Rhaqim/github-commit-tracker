package events

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

var (
	EventChan = make(chan entities.Event)
	ErrorChan = make(chan error, 1)
)

func StartEventListener(processFunc func(event entities.Event) error) {
	go func() {
		for event := range EventChan {
			logger.InfoLogger.Printf("Event received: %v", event)

			if err := processFunc(event); err != nil {
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
