package services

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Rhaqim/savannahtech/internal/api/github"
	"github.com/Rhaqim/savannahtech/internal/api/github/types"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func FetchCommitsByRepository(repoName, pageStr, sizeStr string) ([]entities.Commit, error) {
	size, offset := utils.Paigenation(pageStr, sizeStr)

	return repositories.CommitStore.GetCommitsByRepository(repoName, size, offset)
}

func FetchTopNCommitAuthors(n string) ([]entities.CommitCount, error) {
	nInt, err := strconv.Atoi(n)
	if err != nil {
		return nil, err
	}

	return repositories.CommitStore.GetTopNCommitAuthors(nInt)
}

func ProcessCommitData(owner, repo, start_date string) error {
	var err error

	ownerRepo := fmt.Sprintf("%s/%s", owner, repo)

	url := "https://api.github.com/repos/" + ownerRepo + "/commits"

	if start_date != "" {
		url += "?since=" + start_date
	}

	commitsChan := make(chan []types.Commit)

	go func() {
		err := github.FetchCommits(url, commitsChan)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to fetch commits: %v", err)
		}
		close(commitsChan)
	}()

	for commit := range commitsChan {
		logger.InfoLogger.Println("Received commits: " + strconv.Itoa(len(commit)) + " for " + ownerRepo)

		commits := convertCommitType(commit, ownerRepo)

		if err := repositories.CommitStore.CreateCommits(commits); err != nil {
			return fmt.Errorf("failed to store commits: %w", err)
		}
	}

	// send event to the event queue using go channel

	return err
}

func convertCommitType(commits []types.Commit, ownerRepo string) []entities.Commit {
	var wg sync.WaitGroup

	commitStores := make([]entities.Commit, len(commits))

	for i, commit := range commits {
		wg.Add(1)
		go func(commit types.Commit, i int) {
			defer wg.Done()
			commitStores[i] = entities.Commit{
				SHA:             commit.Sha,
				Author:          commit.Commit.Committer.Name,
				Message:         commit.Commit.Message,
				Date:            commit.Commit.Committer.Date,
				URL:             commit.Commit.Url,
				OwnerRepository: ownerRepo,
			}
		}(commit, i)
	}

	wg.Wait()

	return commitStores

}
