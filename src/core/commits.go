package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"sync"
	"time"

	"github.com/google/uuid"
)

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

func ProcessCommitData(owner, repo string) error {
	log.InfoLogger.Println("Processing commit data")

	var err error
	var commits []types.Commit

	var commitQueue *event.EventQueue = event.NewEventQueue(config.CommitEvent)

	var url string = config.GithubRepoURL + owner + "/" + repo + "/commits"

	commits, err = utils.FetchCommits(url)
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	err = StoreCommit(commits, owner+"/"+repo)
	if err != nil {
		return fmt.Errorf("failed to store commits: %w", err)
	}

	log.InfoLogger.Println("Finished processing commits")

	commitQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Commit data fetched",
		Type:    types.CommitEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}

func ProcessCommitDataChan(owner, repo string) error {
	log.InfoLogger.Println("Processing commit data")

	var err error

	var commitQueue *event.EventQueue = event.NewEventQueue(config.CommitEvent)

	var url string = config.GithubRepoURL + owner + "/" + repo + "/commits"

	var commitsChan = make(chan []types.Commit)

	err = utils.FetchCommitsChan(url, commitsChan)
	if err != nil {
		return fmt.Errorf("failed to fetch commits: %w", err)
	}

	select {
	case commit := <-commitsChan:
		log.InfoLogger.Println("Received commits: ", len(commit))
		err = StoreCommit(commit, owner+"/"+repo)
		if err != nil {
			return fmt.Errorf("failed to store commits: %w", err)
		}
	case <-time.After(time.Hour * 1):
		log.InfoLogger.Println("Timed out")
	}

	log.InfoLogger.Println("Finished processing commits")

	commitQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Commit data fetched",
		Type:    types.CommitEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}
