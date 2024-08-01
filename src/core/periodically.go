package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"time"
)

func PeriodFetch() error {
	log.InfoLogger.Println("Starting periodic commit fetch...")

	var errChan = make(chan error)
	defer close(errChan)

	var commitEvent *event.EventQueue = event.NewEventQueue(config.CommitEvent)

	commitEvent.Subscribe(func(event types.Event) {
		log.InfoLogger.Printf("Commit event received: %v", event)

		err := PeriodicFetch(event.Owner, event.Repo)
		if err != nil {
			errChan <- err
		}
	})

	// Handle errors from the error channel
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process commits: %w", err)
		}
	}

	return nil
}

func PeriodicFetch(owner, repo string) error {
	errChan := make(chan error)
	defer close(errChan)

	var commitStore model.CommitStore

	interval := 1 * time.Hour

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Construct the base URL for fetching commits
	baseURL := config.GithubRepoURL + owner + "/" + repo + "/commits"

	// Start a goroutine to fetch commits periodically
	go func() {
		for range ticker.C {
			// Get the last commit SHA stored
			lastCommitSHA := commitStore.GetLastCommitSHA()

			// Construct the URL with the last commit SHA to fetch new commits
			url := baseURL + "?since=" + lastCommitSHA

			// Fetch the commits from the constructed URL
			commits, err := utils.FetchCommits(url)
			if err != nil {
				log.ErrorLogger.Printf("Error fetching commits: %v", err)
				errChan <- err
				continue
			}

			// Store the fetched commits in the database
			err = StoreCommit(commits)
			if err != nil {
				log.ErrorLogger.Printf("Error storing commits: %v", err)
				errChan <- err
				continue
			}
		}
	}()

	// Handle errors from the error channel
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process commits: %w", err)
		}
	}

	return nil
}
