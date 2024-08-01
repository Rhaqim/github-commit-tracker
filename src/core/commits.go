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
	log.InfoLogger.Println("Processing commit data")

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

func GetEvent() error {
	log.InfoLogger.Println("Starting event listener...")

	var errChan = make(chan error)

	var repoEvent *event.EventQueue = event.NewEventQueue(config.RepoEvent)

	repoEvent.Subscribe(func(event types.Event) {
		log.InfoLogger.Println("Repo event received: ", event)

		// process commit data
		err := ProcessCommitData(event.Owner, event.Repo)
		if err != nil {
			errChan <- err
		}
	})

	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process commits: %w", err)
		}
	}
	return nil
}
