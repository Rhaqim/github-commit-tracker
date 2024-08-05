package services

import (
	"fmt"
	"strconv"
	"strings"
	"sync"

	"github.com/Rhaqim/savannahtech/config"
	"github.com/Rhaqim/savannahtech/internal/api/github"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/core/repositories"
	"github.com/Rhaqim/savannahtech/internal/core/types"
	"github.com/Rhaqim/savannahtech/internal/events"
	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/robfig/cron/v3"
)

func FetchCommitsByRepository(repoName_, pageStr, sizeStr string) ([]entities.Commit, error) {
	size, offset := utils.Paigenation(pageStr, sizeStr)

	repoName := strings.ToLower(repoName_)

	logger.InfoLogger.Printf("Fetching commits for %s\n", repoName)

	return repositories.CommitStore.GetCommitsByRepository(repoName, size, offset)
}

func FetchTopNCommitAuthors(n string) ([]entities.CommitCount, error) {
	nInt, err := strconv.Atoi(n)
	if err != nil {
		return nil, err
	}

	return repositories.CommitStore.GetTopNCommitAuthors(nInt)
}

/*
ProcessCommitData fetches commit data for a repository.

receives the owner, repo, and startDate as arguments.

It fetches the commit data from the GitHub API and stores it in the database.
*/
func ProcessCommitData(owner, repo, startDate string) error {
	var err error

	ownerRepo := fmt.Sprintf("%s/%s", owner, repo)

	url := fmt.Sprintf("https://api.github.com/repos/%s/commits", ownerRepo)

	if startDate != "" {
		url += "?since=" + startDate
	}

	err = processCommit(url, ownerRepo)
	if err != nil {
		return fmt.Errorf("failed to process commit data: %w", err)
	}

	logger.InfoLogger.Printf("Completed initial commit fetching for %s/%s\n", owner, repo)

	event := entities.Event{
		StartDate: startDate,
		Type:      entities.PeriodEvent,
		Owner:     owner,
		Repo:      repo,
	}

	events.SendEvent(event)

	return err
}

/*
PeriodicFetch periodically fetches commit data for a repository.

It uses a cron job to fetch commit data at a specified interval.

It fetches the commit data from the GitHub API and stores it in the database.
*/
func PeriodicFetch(owner, repo string) error {
	interval := config.Config.RefetchInterval

	logger.InfoLogger.Printf("Checking new commits for %s/%s every %s\n", owner, repo, interval)

	c := cron.New()

	ownerRepo := fmt.Sprintf("%s/%s", owner, repo)

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

/*
processCommit fetches commit data from the GitHub API and stores it in the database.
*/
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

/*
convertCommitType converts the commit type from the GitHub API to the Commit type used in the application.
*/
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
