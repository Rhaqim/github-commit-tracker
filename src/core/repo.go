package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"

	"github.com/google/uuid"
)

func ProcessRepositoryData(owner, repo string) error {
	log.InfoLogger.Println("Processing repository data")

	var err error
	var repo_ types.Repository
	var repoStore model.RepositoryStore
	var repoQueue *event.EventQueue = event.NewEventQueue(config.RepoEvent)

	var url string = config.GithubRepoURL + owner + "/" + repo

	repo_, err = utils.FetchRepository(url)
	if err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	repoStore = model.RepositoryStore{
		Name:            repo_.Name,
		Description:     repo_.Description,
		URL:             repo_.URL,
		Language:        repo_.Language,
		StargazersCount: repo_.StargazersCount,
		WatchersCount:   repo_.WatchersCount,
		ForksCount:      repo_.ForksCount,
		RepoCreatedAt:   repo_.CreatedAt,
		RepoUpdatedAt:   repo_.UpdatedAt,
	}

	err = repoStore.InsertRepository()
	if err != nil {
		return fmt.Errorf("failed to insert repository: %w", err)
	}

	repoQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Repository data fetched",
		Type:    types.RepoEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}

func GetRepoEvent() error {
	log.InfoLogger.Println("Starting repo event listener...")

	var errChan = make(chan error)

	var repoEvent *event.EventQueue = event.NewEventQueue(config.NewRepo)

	repoEvent.Subscribe(func(event types.Event) {
		log.InfoLogger.Println("New Repo event received: ", event)

		// process commit data
		err := ProcessRepositoryData(event.Owner, event.Repo)
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
