package core

import (
	"fmt"
	"log"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"

	"github.com/google/uuid"
)

func ProcessRepositoryData(owner, repo string) error {
	log.Println("Processing repository data")

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
