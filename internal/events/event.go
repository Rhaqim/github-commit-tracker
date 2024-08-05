package events

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

var (
	EventChan = make(chan entities.Event, 4)
	ErrorChan = make(chan error, 1)
)

/*
StartEventListener starts the event listener that listens for events on the EventChan and processes them using the provided processFunc.
*/
func StartEventListener(processFunc func(event entities.Event) error) {
	go func() {
		for event := range EventChan {
			logger.InfoLogger.Printf("Event received: %v", event)

			go func(e entities.Event) {
				if err := processFunc(e); err != nil {
					ErrorChan <- err
				}
			}(event)
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

/*
SendEvent sends an event to the EventChan.
*/
func SendEvent(event entities.Event) {
	EventChan <- event
}
