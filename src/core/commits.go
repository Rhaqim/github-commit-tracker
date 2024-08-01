package core

import (
	"fmt"
	"log"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"sync"

	"github.com/google/uuid"
)

func StoreCommit(commits []types.Commit) error {
	var wg sync.WaitGroup
	var commitStore model.CommitStore

	commitStores := make([]model.CommitStore, len(commits))

	for i, commit := range commits {
		wg.Add(1)
		go func(commit types.Commit, i int) {
			defer wg.Done()
			commitStores[i] = model.CommitStore{
				SHA:     commit.Sha,
				Author:  commit.Commit.Committer.Name,
				Message: commit.Commit.Message,
				Date:    commit.Commit.Committer.Date,
				URL:     commit.Commit.Url,
			}
		}(commit, i)
	}

	wg.Wait()

	err := commitStore.InsertManyCommits(commitStores)
	if err != nil {
		return fmt.Errorf("failed to insert commits: %w", err)
	}

	return nil
}

func ProcessCommitData(owner, repo string) error {
	log.Println("Processing commit data")

	var err error
	var commits []types.Commit

	var commitQueue *event.EventQueue = event.NewEventQueue(config.CommitEvent)

	var url string = config.GithubRepoURL + owner + "/" + repo + "/commits"

	commits, err = utils.FetchCommits(url)
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	err = StoreCommit(commits)
	if err != nil {
		return fmt.Errorf("failed to store commits: %w", err)
	}

	log.Println("Finished processing commits")

	commitQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Commit data fetched",
		Type:    types.CommitEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}

func GetEvent() error {
	log.Println("Starting event listener...")

	var errChan = make(chan error)

	var repoEvent *event.EventQueue = event.NewEventQueue(config.RepoEvent)

	repoEvent.Subscribe(func(event types.Event) {
		log.Println("Repo event received: ", event)

		// process commit data
		err := ProcessCommitData(event.Owner, event.Repo)
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

// Generalized function to fetch and process data
// func FetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc[T]) ([]T, error), processFunc func(T) error) error {
// 	var wg sync.WaitGroup
// 	var commitQueue *event.EventQueue = event.NewEventQueue("commit-event")

// 	var data []T

// 	data, err := fetchFunc(owner, repo, utils.MakeRequest2[T])
// 	if err != nil {
// 		return fmt.Errorf("failed to fetch data: %w", err)
// 	}

// 	var errChan = make(chan error)

// 	for _, item := range data {
// 		wg.Add(1)
// 		go func(item T) {
// 			defer wg.Done()
// 			if err := processFunc(item); err != nil {
// 				errChan <- err
// 			}
// 		}(item)
// 	}

// 	wg.Wait()
// 	close(errChan)

// 	for err := range errChan {
// 		if err != nil {
// 			return fmt.Errorf("failed to process data: %w", err)
// 		}
// 	}

// 	log.Println("Finished processing commits")

// 	commitQueue.Publish(types.Event{
// 		ID:      uuid.New().String(),
// 		Message: "Commit data fetched",
// 		Type:    types.CommitEvent,
// 		Owner:   owner,
// 		Repo:    repo,
// 	})

// 	return nil
// }

// func FetchCommit(owner, repo string, makeRequest types.RequestFunc[types.Commit]) ([]types.Commit, error) {
// 	url := "https://api.github.com/repos/" + owner + "/" + repo + "/commits"
// 	return utils.FetchData[types.Commit](url, makeRequest)
// }

// // Specific process functions for commits and repositories
// func processCommit(commit types.Commit) error {
// 	commitStore := model.CommitStore{
// 		SHA:     commit.Sha,
// 		Author:  commit.Commit.Committer.Name,
// 		Message: commit.Commit.Message,
// 		Date:    commit.Commit.Committer.Date,
// 		URL:     commit.Commit.Url,
// 	}

// 	return commitStore.InsertCommit()
// }

// // Simplified functions for commits and repositories
// func CommitData(owner, repo string) error {
// 	log.Println("Fetching commits for", owner, repo)
// 	// subscribe to event
// 	return FetchAndProcessData[types.Commit](owner, repo, FetchCommit, processCommit)
// }
