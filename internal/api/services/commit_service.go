package services

import (
	"fmt"
	"strconv"
	"sync"

	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/api/github"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/core/types"
	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/robfig/cron/v3"
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

func ProcessCommitData(event_ entities.Event) error {
	var err error

	var owner, repo, start_date string = event_.Owner, event_.Repo, event_.StartDate

	if event_.Type != entities.CommitEvent {
		return PeriodicFetch(owner, repo, start_date)
	}

	ownerRepo := fmt.Sprintf("%s/%s", owner, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits", ownerRepo)

	if start_date != "" {
		url += "?since=" + start_date
	}

	err = processCommit(url, ownerRepo)
	if err != nil {
		return fmt.Errorf("failed to process commit data: %w", err)
	}

	logger.InfoLogger.Printf("Completed initial commit fetching for %s/%s\n", owner, repo)

	err = PeriodicFetch(owner, repo, start_date)
	if err != nil {
		return fmt.Errorf("failed to start periodic commit fetching: %w", err)
	}

	return err
}

/*
PeriodicFetch periodically fetches commit data for a repository.

It uses a cron job to fetch commit data at a specified interval.

It fetches the commit data from the GitHub API and stores it in the database.
*/
func PeriodicFetch(owner, repo, _ string) error {
	interval := config.Config.RefetchInterval

	logger.InfoLogger.Printf("Started periodic commit fetching for %s/%s every %s\n", owner, repo, interval)

	c := cron.New()

	ownerRepo := owner + "/" + repo

	// Construct the base URL for fetching commits
	baseURL := fmt.Sprintf("https://api.github.com/repos/%s/commits", ownerRepo)

	c.AddFunc(fmt.Sprintf("@every %s", interval), func() {
		// Get the last commit SHA stored
		lastCommitDate := repositories.CommitStore.GetLastCommitDate(ownerRepo)

		// Construct the URL with the last commit date to fetch new commits
		url := fmt.Sprintf("%s?since=%s", baseURL, lastCommitDate)

		if lastCommitDate == "" {
			url = baseURL
		}

		err := processCommit(url, ownerRepo)
		if err != nil {
			logger.ErrorLogger.Printf("Failed to process commit data: %v", err)
			return
		}

	})
	c.Start()

	return nil
}

func processCommit(url, ownerRepo string) error {
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

		if len(commit) == 0 {
			break
		}

		commits := convertCommitType(commit, ownerRepo)

		if err := repositories.CommitStore.CreateCommits(commits); err != nil {
			return fmt.Errorf("failed to store commits: %w", err)
		}
	}

	err := repositories.RepoStore.UpdateRepositoryIndexed(ownerRepo)
	if err != nil {
		return fmt.Errorf("failed to update repository: %w", err)
	}

	return nil
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
