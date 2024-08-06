package events

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

var (
	EventChan = make(chan entities.Event, 4)
)

/*
StartEventListener starts the event listener that listens for events on the EventChan and processes them using the provided processFunc.
*/
func StartEventListener(processFunc func(event entities.Event)) {
	go func() {
		for event := range EventChan {
			logger.InfoLogger.Printf("Event received: %v", event)

			go func(e entities.Event) {
				processFunc(e)
			}(event)
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
