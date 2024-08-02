package core

import (
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/types"
)

/*
startEventListener starts a new event listener for a given event key and process function.

It creates a new event queue with the given event key and subscribes to it. It then starts a goroutine to handle errors from the subscription.

The goroutine logs any errors that occur during the subscription and handles them appropriately.

It then starts a goroutine to handle events from the event queue. It logs any errors that occur during the processing of events and handles them appropriately.

Finally, it returns nil if no errors occurred during the startup of the event listener.
*/
func startEventListener(eventKey string, processFunc func(owner, repo string) error, listenerName string) error {
	log.InfoLogger.Printf("Starting %s listener...", listenerName)

	errChan := make(chan error, 1)
	defer close(errChan)

	eventQueue := event.NewEventQueue(eventKey)

	// Start a goroutine to handle errors from the subscription
	go func() {
		for err := range errChan {
			if err != nil {
				// Log and handle the error appropriately
				log.ErrorLogger.Printf("Error processing event in %s: %v", listenerName, err)
			}
		}
	}()

	eventQueue.Subscribe(func(event types.Event) {
		log.InfoLogger.Printf("%s event received: %v", listenerName, event)

		if err := processFunc(event.Owner, event.Repo); err != nil {
			errChan <- err
		}
	})

	return nil
}

func GetRepoEvent() error {
	return startEventListener(config.NewRepo, ProcessRepositoryData, "repo event")
}

func GetCommitEvent() error {
	return startEventListener(config.RepoEvent, ProcessCommitData, "commit event")
}

func PeriodFetch() error {
	return startEventListener(config.CommitEvent, PeriodicFetch, "periodic commit fetch")
}
