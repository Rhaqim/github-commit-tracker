package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/event"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

/*
ProcessRepositoryData processes the repository data for a repository.

It fetches the repository data from the GitHub API and stores it in the database.

It also publishes an event to the event queue indicating that the repository data has been fetched.
*/
func ProcessRepositoryData(owner, repo, _ string) error {
	log.InfoLogger.Println("Processing repository data")

	repoStore := model.RepositoryStore{}

	repoQueue := event.NewEventQueue(config.CommitEvent)
	commitQueue := event.NewEventQueue(config.PeriodEvent)

	ownerRepo := owner + "/" + repo

	url := config.GithubRepoURL + ownerRepo

	// Check if the repository exists in the store
	err := repoStore.GetRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handleNewRepository(url, owner, repo, &repoStore, repoQueue, commitQueue)
		}
		return fmt.Errorf("failed to get repository: %w", err)
	}

	return handleExistingRepository(owner, repo, &repoStore, commitQueue)
}

func handleNewRepository(url, owner, repo string, repoStore *model.RepositoryStore, repoQueue, commitQueue *event.EventQueue) error {
	log.InfoLogger.Println("Handling new repository for " + owner + "/" + repo)

	// Fetch repository data from the remote source
	repo_, err := utils.FetchRepository(url)
	if err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	// Populate repository store with fetched data
	*repoStore = model.RepositoryStore{
		Name:            repo_.Name,
		Description:     repo_.Description,
		URL:             repo_.URL,
		Language:        repo_.Language,
		StargazersCount: repo_.StargazersCount,
		WatchersCount:   repo_.WatchersCount,
		ForksCount:      repo_.ForksCount,
		RepoCreatedAt:   repo_.CreatedAt,
		RepoUpdatedAt:   repo_.UpdatedAt,
		OwnerRepository: owner + "/" + repo,
	}

	// Insert the repository into the store
	err = repoStore.InsertRepository()
	if err != nil {
		if strings.Contains(err.Error(), "unique constraint") {
			return handleExistingRepository(owner, repo, repoStore, commitQueue)
		}
		return fmt.Errorf("failed to insert repository: %w", err)
	}

	// Publish new repository event
	repoQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "Repository data fetched",
		Type:    types.RepoEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}

func handleExistingRepository(owner, repo string, repoStore *model.RepositoryStore, commitQueue *event.EventQueue) error {
	log.InfoLogger.Println("Handling existing repository for " + owner + "/" + repo)

	// Check if the repository is already indexed
	err := repoStore.GetRepositoryByOwnerRepo(owner + "/" + repo)
	if err != nil {
		return fmt.Errorf("failed to get repository: %w", err)
	}

	if repoStore.Indexed {
		log.InfoLogger.Println("Repository already indexed, sending commit data")

		// Publish commit event for periodic fetching of commits
		commitQueue.Publish(types.Event{
			ID:      uuid.New().String(),
			Message: "Commit data fetched",
			Type:    types.CommitEvent,
			Owner:   owner,
			Repo:    repo,
		})
		return nil
	}

	// Publish new repository event if not indexed
	commitQueue.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "New repository event",
		Type:    types.NewRepo,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}

func LoadStartupRepo() error {
	// wait for 2 seconds to allow the event listeners to start
	<-time.After(2 * time.Second)

	log.InfoLogger.Println("Loading startup repository")

	var newRepoEvent *event.EventQueue = event.NewEventQueue(config.RepoEvent)
	if err := newRepoEvent.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "New repository event",
		Type:    types.NewRepo,
		Owner:   strings.ToLower(config.DefaultOwner),
		Repo:    strings.ToLower(config.DefaultRepo),
	}); err != nil {
		return fmt.Errorf("failed to publish startup repository event: %w", err)
	}

	return nil
}
