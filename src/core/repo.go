package core

import (
	"encoding/json"
	"fmt"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
)

func FetchRepo(owner, repo string, makeRequest types.RequestFunc[byte]) ([]byte, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo
	return utils.FetchData(url, makeRequest)
}

func processRepository(repo types.Repository) error {
	repoStore := model.RepositoryStore{
		Name:            repo.Name,
		Description:     repo.Description,
		URL:             repo.URL,
		Language:        repo.Language,
		StargazersCount: repo.StargazersCount,
		WatchersCount:   repo.WatchersCount,
		ForksCount:      repo.ForksCount,
		RepoCreatedAt:   repo.CreatedAt,
		RepoUpdatedAt:   repo.UpdatedAt,
	}
	return repoStore.InsertRepository()
}

func fetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc[byte]) ([]byte, error), processFunc func(T) error) error {

	var data T

	body, err := fetchFunc(owner, repo, utils.MakeRequest)
	if err != nil {
		return fmt.Errorf("fetch error: %v", err)
	}

	err = json.Unmarshal(body, &data)
	if err != nil {
		return fmt.Errorf("unmarshal error: %v", err)
	}

	err = processFunc(data)
	if err != nil {
		return fmt.Errorf("process error: %v", err)
	}

	return nil
}

func RepositoryData(owner, repo string) error {
	err := fetchAndProcessData[types.Repository](owner, repo, FetchRepo, processRepository)
	if err != nil {
		return err
	}

	// publish event to redis
	var repoEvent *event.EventQueue = event.NewEventQueue("repo-event")
	repoEvent.Publish(types.Event{
		Message: "Repository data fetched",
		Type:    types.RepoEvent,
		Owner:   owner,
		Repo:    repo,
	})

	return nil
}
