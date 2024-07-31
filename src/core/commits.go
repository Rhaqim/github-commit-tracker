package core

import (
	"fmt"
	"log"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"sync"
)

// Generalized function to fetch and process data
func FetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc[T]) ([]T, error), processFunc func(T) error) error {
	var wg sync.WaitGroup

	var data []T

	data, err := fetchFunc(owner, repo, utils.MakeRequest2[T])
	if err != nil {
		return fmt.Errorf("failed to fetch data: %w", err)
	}

	var errChan = make(chan error)

	for _, item := range data {
		wg.Add(1)
		go func(item T) {
			defer wg.Done()
			if err := processFunc(item); err != nil {
				errChan <- err
			}
		}(item)
	}

	wg.Wait()
	close(errChan)

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process data: %w", err)
		}
	}

	return nil
}

func FetchCommit(owner, repo string, makeRequest types.RequestFunc[types.Commit]) ([]types.Commit, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/commits"
	return utils.FetchData[types.Commit](url, makeRequest)
}

// Specific process functions for commits and repositories
func processCommit(commit types.Commit) error {
	commitStore := model.CommitStore{
		SHA:     commit.Sha,
		Author:  commit.Commit.Committer.Name,
		Message: commit.Commit.Message,
		Date:    commit.Commit.Committer.Date,
		URL:     commit.Commit.Url,
	}

	return commitStore.InsertCommit()
}

// Simplified functions for commits and repositories
func CommitData(owner, repo string) error {
	log.Println("Fetching commits for", owner, repo)
	// subscribe to event
	return FetchAndProcessData[types.Commit](owner, repo, FetchCommit, processCommit)
}

func GetEvent() error {
	log.Println("Starting event listener...")

	var errChan = make(chan error)

	var commitEvent *event.EventQueue = event.NewEventQueue("repo-event")

	commitEvent.Subscribe(func(event types.Event) {
		log.Println("Got event: ", event)
		// process commit data
		err := CommitData(event.Owner, event.Repo)
		if err != nil {
			errChan <- err
		}
	})

	for err := range errChan {
		if err != nil {
			return err
		}
	}
	return nil
}
