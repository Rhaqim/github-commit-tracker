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

// func FetchRepo(owner, repo string, makeRequest types.RequestFunc[byte]) ([]byte, error) {
// 	url := "https://api.github.com/repos/" + owner + "/" + repo
// 	return utils.FetchData(url, makeRequest)
// }

// func processRepository(repo types.Repository) error {
// 	repoStore := model.RepositoryStore{
// 		Name:            repo.Name,
// 		Description:     repo.Description,
// 		URL:             repo.URL,
// 		Language:        repo.Language,
// 		StargazersCount: repo.StargazersCount,
// 		WatchersCount:   repo.WatchersCount,
// 		ForksCount:      repo.ForksCount,
// 		RepoCreatedAt:   repo.CreatedAt,
// 		RepoUpdatedAt:   repo.UpdatedAt,
// 	}
// 	return repoStore.InsertRepository()
// }

// func fetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc[byte]) ([]byte, error), processFunc func(T) error) error {

// 	var data T

// 	body, err := fetchFunc(owner, repo, utils.MakeRequest)
// 	if err != nil {
// 		return fmt.Errorf("fetch error: %v", err)
// 	}

// 	err = json.Unmarshal(body, &data)
// 	if err != nil {
// 		return fmt.Errorf("unmarshal error: %v", err)
// 	}

// 	err = processFunc(data)
// 	if err != nil {
// 		return fmt.Errorf("process error: %v", err)
// 	}

// 	return nil
// }

// func RepositoryData(owner, repo string) error {
// 	err := fetchAndProcessData[types.Repository](owner, repo, FetchRepo, processRepository)
// 	if err != nil {
// 		return err
// 	}

// 	// publish event to redis
// 	var repoEvent *event.EventQueue = event.NewEventQueue("repo-event")
// 	repoEvent.Publish(types.Event{
// 		Message: "Repository data fetched",
// 		Type:    types.RepoEvent,
// 		Owner:   owner,
// 		Repo:    repo,
// 	})

// 	return nil
// }
