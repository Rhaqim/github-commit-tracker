package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/types"
)

func startEventListener(eventKey string, processFunc func(owner, repo string) error, listenerName string) error {
	log.InfoLogger.Printf("Starting %s listener...", listenerName)

	errChan := make(chan error)
	defer close(errChan)

	eventQueue := event.NewEventQueue(eventKey)

	eventQueue.Subscribe(func(event types.Event) {
		log.InfoLogger.Printf("%s event received: %v", listenerName, event)

		if err := processFunc(event.Owner, event.Repo); err != nil {
			errChan <- err
		}
	})

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process event in %s: %w", listenerName, err)
		}
	}

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
