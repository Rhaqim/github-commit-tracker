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

	"github.com/google/uuid"
)

func ProcessRepositoryData(owner, repo string) error {
	log.InfoLogger.Println("Processing repository data")

	repoStore := model.RepositoryStore{}
	repoQueue := event.NewEventQueue(config.RepoEvent)
	commitQueue := event.NewEventQueue(config.CommitEvent)
	ownerRepo := owner + "/" + repo
	url := config.GithubRepoURL + ownerRepo

	// Check if the repository exists in the store
	err := repoStore.GetRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
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

// func ProcessRepositoryData(owner, repo string) error {
// 	log.InfoLogger.Println("Processing repository data")

// 	var err error
// 	var repo_ types.Repository
// 	var repoStore model.RepositoryStore

// 	var repoQueue *event.EventQueue = event.NewEventQueue(config.RepoEvent)
// 	var commitQueue *event.EventQueue = event.NewEventQueue(config.CommitEvent)

// 	var url string = config.GithubRepoURL + owner + "/" + repo

// 	// check if the repository exists
// 	if err = repoStore.GetRepositoryByOwnerRepo(owner + "/" + repo); err != nil {
// 		if err.Error() == "sql: no rows in result set" {

// 			repo_, err = utils.FetchRepository(url)
// 			if err != nil {
// 				return fmt.Errorf("failed to fetch repository: %w", err)
// 			}

// 			repoStore = model.RepositoryStore{
// 				Name:            repo_.Name,
// 				Description:     repo_.Description,
// 				URL:             repo_.URL,
// 				Language:        repo_.Language,
// 				StargazersCount: repo_.StargazersCount,
// 				WatchersCount:   repo_.WatchersCount,
// 				ForksCount:      repo_.ForksCount,
// 				RepoCreatedAt:   repo_.CreatedAt,
// 				RepoUpdatedAt:   repo_.UpdatedAt,
// 				OwnerRepository: owner + "/" + repo,
// 			}

// 			err = repoStore.InsertRepository()
// 			if err != nil {
// 				// Check if the error is a unique constraint violation
// 				if strings.Contains(err.Error(), "unique constraint") {
// 					// get the repository from the database and check if it's indexed
// 					var repository model.RepositoryStore
// 					err = repository.GetRepositoryByOwnerRepo(owner + "/" + repo)
// 					if err != nil {
// 						return fmt.Errorf("failed to get repository: %w", err)
// 					}

// 					if repository.Indexed {
// 						// return fmt.Errorf("repository already indexed, skipping")

// 						log.InfoLogger.Println("Repository already intially indexed, sending commit data")

// 						// if already indexed, publish to the periodic repo event  to continue fetching commits
// 						commitQueue.Publish(types.Event{
// 							ID:      uuid.New().String(),
// 							Message: "Commit data fetched",
// 							Type:    types.CommitEvent,
// 							Owner:   owner,
// 							Repo:    repo,
// 						})

// 						return nil
// 					}

// 					// publish a new repository event
// 					repoQueue.Publish(types.Event{
// 						ID:      uuid.New().String(),
// 						Message: "New repository event",
// 						Type:    types.NewRepo,
// 						Owner:   owner,
// 						Repo:    repo,
// 					})
// 				}
// 			}

// 			repoQueue.Publish(types.Event{
// 				ID:      uuid.New().String(),
// 				Message: "Repository data fetched",
// 				Type:    types.RepoEvent,
// 				Owner:   owner,
// 				Repo:    repo,
// 			})
// 		}

// 		return fmt.Errorf("failed to get repository: %w", err)
// 	}

// 	return nil
// }

func LoadStartupRepo() error {
	var newRepoEvent *event.EventQueue = event.NewEventQueue(config.NewRepo)
	if err := newRepoEvent.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "New repository event",
		Type:    types.NewRepo,
		Owner:   config.DefaultOwner,
		Repo:    config.DefaultRepo,
	}); err != nil {
		return fmt.Errorf("failed to publish startup repository event: %w", err)
	}

	return nil
}
