package services

import (
	"fmt"
	"strings"
	"time"

	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/api/github"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/core/types"
	"github.com/Rhaqim/savannahtech/internal/events"
	"github.com/Rhaqim/savannahtech/pkg/logger"
	"gorm.io/gorm"
)

func FetchRepositoryByOwnerRepo(owner, repo string) (entities.Repository, error) {
	ownerRepo := strings.ToLower(fmt.Sprintf("%s/%s", owner, repo))
	return repositories.RepoStore.GetRepositoryByOwnerRepo(ownerRepo)
}

/*
ProcessRepositoryData processes the repository data for a repository.

It fetches the repository data from the GitHub API and stores it in the database.

It also publishes an event to the event queue indicating that the repository data has been fetched.
*/
func ProcessRepository(owner, repo, startDate string) error {

	url := fmt.Sprintf("https://api.github.com/repos/%s/%s", owner, repo)

	// Check if the repository exists in the store
	repo_, err := FetchRepositoryByOwnerRepo(owner, repo)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return handleNewRepository(url, owner, repo, startDate)
		}
		return err
	}

	return handleExistingRepository(repo_, owner, repo, startDate)
}

/*
handleNewRepository handles the case where the repository is new.

It fetches the repository data from the GitHub API and stores it in the database.

It also publishes an event to the event queue indicating that the repository data has been fetched.
*/
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
		StartDate: startDate,
		Type:      entities.CommitEvent,
		Owner:     owner,
		Repo:      repo,
	}

	events.SendEvent(event)

	return nil
}

/*
handleExistingRepository handles the case where the repository already exists.

It checks if the repository is indexed and sends the appropriate event to the event queue.
*/
func handleExistingRepository(repo_ entities.Repository, owner, repo, startDate string) error {
	logger.InfoLogger.Printf("Handling existing repository %s/%s\n", owner, repo)

	if repo_.Indexed {
		logger.InfoLogger.Println("Repository already indexed, sending commit data")

		// Publish commit event for periodic fetching of commits
		event := entities.Event{
			StartDate: startDate,
			Type:      entities.PeriodEvent,
			Owner:     owner,
			Repo:      repo,
		}

		events.SendEvent(event)

		return nil
	}

	// Publish new repository event if not indexed
	event := entities.Event{
		StartDate: startDate,
		Type:      entities.CommitEvent,
		Owner:     owner,
		Repo:      repo,
	}

	events.SendEvent(event)

	return nil
}

/*
convertRepositoryType converts the repository type from the API to the internal repository type.
*/
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

/*
LoadStartupRepo loads the startup repository.

the configuration from the .env file is used to determine the repository to load.

It sends an event to the event queue to fetch the repository data.
*/
func LoadStartupRepo() {
	// wait for 2 seconds to allow the event listeners to start
	<-time.After(2 * time.Second)

	owner := config.Config.DefaultOwner
	repo := config.Config.DefaultRepo
	startDate := config.Config.DefaultStartDate

	logger.InfoLogger.Printf("Loading startup repository %s/%s\n", owner, repo)

	event := entities.Event{
		StartDate: startDate,
		Type:      entities.RepoEvent,
		Owner:     owner,
		Repo:      repo,
	}

	events.SendEvent(event)

}
