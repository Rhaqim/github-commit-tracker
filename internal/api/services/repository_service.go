package services

import (
	"fmt"

	"github.com/Rhaqim/savannahtech/internal/api/github"
	"github.com/Rhaqim/savannahtech/internal/api/github/types"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/events"
	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/google/uuid"
	"gorm.io/gorm"
)

func FetchRepositoryByOwnerRepo(ownerRepo string) (entities.Repository, error) {
	return repositories.RepoStore.GetRepositoryByOwnerRepo(ownerRepo)
}

/*
ProcessRepositoryData processes the repository data for a repository.

It fetches the repository data from the GitHub API and stores it in the database.

It also publishes an event to the event queue indicating that the repository data has been fetched.
*/
func ProcessRepository(owner, repo, startDate string) error {

	ownerRepo := fmt.Sprintf("%s/%s", owner, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s", ownerRepo)

	// Check if the repository exists in the store
	repo_, err := FetchRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handleNewRepository(url, owner, repo, startDate)
		}
		return err
	}

	return handleExistingRepository(repo_, owner, repo)
}

func convertRepositoryType(repo_ types.Repository, ownerRepo string) entities.Repository {
	var repo entities.Repository = entities.Repository{
		Name:            repo_.Name,
		Description:     repo_.Description,
		URL:             repo_.URL,
		Language:        repo_.Language,
		StargazersCount: repo_.StargazersCount,
		WatchersCount:   repo_.WatchersCount,
		ForksCount:      repo_.ForksCount,
		RepoCreatedAt:   repo_.CreatedAt,
		RepoUpdatedAt:   repo_.UpdatedAt,
		OwnerRepository: ownerRepo,
	}

	return repo
}

func handleNewRepository(url, owner, repo, startDate string) error {
	logger.InfoLogger.Printf("Handling new repository %s/%s\n", owner, repo)

	// Fetch repository data from the remote source
	repo_, err := github.FetchRepository(url)
	if err != nil {
		return fmt.Errorf("failed to fetch repository: %w", err)
	}

	// Populate repository store with fetched data
	err = repositories.RepoStore.StoreRepository(convertRepositoryType(repo_, fmt.Sprintf("%s/%s", owner, repo)))
	if err != nil {
		return fmt.Errorf("failed to insert repository: %w", err)
	}

	logger.InfoLogger.Printf("Successfully accessed repository %s/%s\n", owner, repo)

	event := entities.Event{
		ID:      uuid.New().String(),
		From:    startDate,
		Message: "New repository event",
		Type:    entities.NewRepo,
		Owner:   owner,
		Repo:    repo,
	}

	events.SendEvent(event)

	return nil
}

func handleExistingRepository(repo_ entities.Repository, owner, repo string) error {
	logger.InfoLogger.Println("Handling existing repository for")

	if repo_.Indexed {
		logger.InfoLogger.Println("Repository already indexed, sending commit data")

		// Publish commit event for periodic fetching of commits
		event := entities.Event{
			ID:      uuid.New().String(),
			From:    "",
			Message: "Commit data fetched",
			Type:    entities.CommitEvent,
			Owner:   "",
			Repo:    "",
		}

		events.SendEvent(event)

		return nil
	}

	// Publish new repository event if not indexed
	event := entities.Event{
		ID:      uuid.New().String(),
		From:    "startDate",
		Message: "New repository event",
		Type:    entities.NewRepo,
		Owner:   owner,
		Repo:    repo,
	}

	events.SendEvent(event)

	return nil
}
