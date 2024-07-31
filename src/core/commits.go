package core

import (
	"encoding/json"
	"fmt"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"sync"
)

// Generalized function to fetch and process data
func FetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc) ([]byte, error), processFunc func(T) error) error {
	var wg sync.WaitGroup

	var data []T

	body, err := fetchFunc(owner, repo, utils.MakeRequest)
	if err != nil {
		return err
	}

	fmt.Println(string(body))

	err = json.Unmarshal(body, &data)
	if err != nil {
		return err
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
			return err
		}
	}

	return nil
}

func FetchCommit(owner, repo string, makeRequest types.RequestFunc) ([]byte, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo + "/commits"
	return utils.FetchData(url, makeRequest)
}

// Specific process functions for commits and repositories
func processCommit(commit types.Commit) error {
	commitStore := model.CommitStore{
		Author:  commit.Commit.Committer.Name,
		Message: commit.Commit.Message,
		Date:    commit.Commit.Committer.Date,
		URL:     commit.Commit.Url,
	}
	fmt.Println("Storing commit", commitStore)

	return commitStore.InsertCommit()
}

// Simplified functions for commits and repositories
func CommitData(owner, repo string) error {
	fmt.Println("Fetching commits for", owner, repo)
	// subscribe to event
	return FetchAndProcessData[types.Commit](owner, repo, FetchCommit, processCommit)
}

func GetEvent() error {
	fmt.Println("Starting event listener...")

	var errChan = make(chan error)

	var commitEvent *event.EventQueue = event.NewEventQueue("repo-event")

	commitEvent.Subscribe(func(event event.Event) {
		fmt.Println("Got event: ", event)
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
