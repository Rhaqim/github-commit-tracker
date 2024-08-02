package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"strconv"
	"sync"

	"github.com/google/uuid"
)

/*
StoreCommit stores the commit data in the database.
*/
func StoreCommit(commits []types.Commit, ownerRepo string) error {
	var wg sync.WaitGroup
	var commitStore model.CommitStore

	commitStores := make([]model.CommitStore, len(commits))

	for i, commit := range commits {
		wg.Add(1)
		go func(commit types.Commit, i int) {
			defer wg.Done()
			commitStores[i] = model.CommitStore{
				SHA:             commit.Sha,
				Author:          commit.Commit.Committer.Name,
				Message:         commit.Commit.Message,
				Date:            commit.Commit.Committer.Date,
				URL:             commit.Commit.Url,
				OwnerRepository: ownerRepo,
			}
		}(commit, i)
	}

	wg.Wait()

	err := commitStore.InsertManyCommits(commitStores)
	if err != nil {
		return fmt.Errorf("failed to insert commits: %w", err)
	}

	// update the repository to indicate that it has been indexed
	var repository model.RepositoryStore
	err = repository.GetRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	repository.Indexed = true
	err = repository.UpdateRepository()
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	return nil
}

/*
ProcessCommitData processes the commit data for a repository.

It fetches the commit data from the GitHub API and stores it in the database.

It also publishes an event to the event queue indicating that the commit data has been fetched.
*/
func ProcessCommitData(owner, repo, fromDate string) error {
	log.InfoLogger.Println("Processing commit data")

	commitQueue := event.NewEventQueue(config.PeriodEvent)

	ownerRepo := owner + "/" + repo

	url := config.GithubRepoURL + ownerRepo + "/commits"

	if fromDate != "" {
		url += "?since=" + fromDate
	}

	commitsChan := make(chan []types.Commit)

	// Fetch commits in a separate goroutine
	go func() {
		err := utils.FetchCommits(url, commitsChan)
		if err != nil {
			log.ErrorLogger.Printf("Failed to fetch commits: %v", err)
		}
		close(commitsChan) // Close channel after fetching all commits
	}()

	for commit := range commitsChan {
		log.InfoLogger.Println("Received commits: " + strconv.Itoa(len(commit)) + " for " + ownerRepo)

		if err := StoreCommit(commit, ownerRepo); err != nil {
			return fmt.Errorf("failed to store commits: %w", err)
		}
	}

	log.InfoLogger.Println("Finished indexing initial commits")

	commitQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Commit data fetched",
		Type:    types.CommitEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}
