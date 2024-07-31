package core

import (
	"encoding/json"
	"fmt"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
)

func FetchRepo(owner, repo string, makeRequest types.RequestFunc) ([]byte, error) {
	url := "https://api.github.com/repos/" + owner + "/" + repo
	return utils.FetchData(url, makeRequest)
}

func processRepository(repo model.RepositoryStore) error {

	return repo.InsertRepository()
}

func fetchAndProcessData[T any](owner, repo string, fetchFunc func(string, string, types.RequestFunc) ([]byte, error), processFunc func(T) error) error {

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
	err := fetchAndProcessData[model.RepositoryStore](owner, repo, FetchRepo, processRepository)
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
