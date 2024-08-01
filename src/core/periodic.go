package core

import (
	"fmt"
	"savannahtech/src/config"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/utils"
	"time"
)

func PeriodicFetch(owner, repo string) error {
	errChan := make(chan error)
	defer close(errChan)

	var commitStore model.CommitStore

	interval := 1 * time.Hour

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	// Construct the base URL for fetching commits
	baseURL := config.GithubRepoURL + owner + "/" + repo + "/commits"

	// Start a goroutine to fetch commits periodically
	go func() {
		for range ticker.C {
			// Get the last commit SHA stored
			lastCommitDate := commitStore.GetLastCommitDate()

			// Construct the URL with the last commit SHA to fetch new commits
			url := baseURL + "?since=" + lastCommitDate.String()

			// Fetch the commits from the constructed URL
			commits, err := utils.FetchCommits(url)
			if err != nil {
				log.ErrorLogger.Printf("Error fetching commits: %v", err)
				errChan <- err
				continue
			}

			// Store the fetched commits in the database
			err = StoreCommit(commits, owner+"/"+repo)
			if err != nil {
				log.ErrorLogger.Printf("Error storing commits: %v", err)
				errChan <- err
				continue
			}
		}
	}()

	// Handle errors from the error channel
	for err := range errChan {
		if err != nil {
			return fmt.Errorf("failed to process commits: %w", err)
		}
	}

	return nil
}
