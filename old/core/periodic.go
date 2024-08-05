package core

import (
	"strconv"

	"github.com/Rhaqim/savannahtech/old/config"
	"github.com/Rhaqim/savannahtech/old/model"
	"github.com/Rhaqim/savannahtech/old/types"
	"github.com/Rhaqim/savannahtech/old/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"

	"github.com/robfig/cron/v3"
)

/*
PeriodicFetch periodically fetches commit data for a repository.

It uses a cron job to fetch commit data at a specified interval.

It fetches the commit data from the GitHub API and stores it in the database.
*/
func PeriodicFetch(owner, repo, _ string) error {
	logger.InfoLogger.Printf("Started periodic commit fetching for %s/%s every %s\n", owner, repo, config.RefetchInterval)

	c := cron.New()
	var commitStore model.CommitStore

	ownerRepo := owner + "/" + repo

	// Construct the base URL for fetching commits
	baseURL := config.GithubRepoURL + ownerRepo + "/commits"

	c.AddFunc("@every "+config.RefetchInterval, func() {
		// Get the last commit SHA stored
		lastCommitDate := commitStore.GetLastCommitDate(ownerRepo)

		// Construct the URL with the last commit SHA to fetch new commits
		url := baseURL + "?since=" + lastCommitDate.Format("2006-01-02T15:04:05Z")

		commitsChan := make(chan []types.Commit)

		// Fetch commits in a separate goroutine
		go func() {
			err := utils.FetchCommits(url, commitsChan)
			if err != nil {
				logger.ErrorLogger.Printf("Failed to fetch commits: %v", err)
				return
			}
			close(commitsChan) // Close channel after fetching all commits
		}()

		for commit := range commitsChan {
			logger.InfoLogger.Println("Received commits: " + strconv.Itoa(len(commit)) + " for " + owner + "/" + repo)

			// check if the commit is empty
			if len(commit) == 0 {
				continue
			}

			if err := StoreCommit(commit, owner+"/"+repo); err != nil {

				logger.ErrorLogger.Printf("Failed to store commits: %v", err)
				return
			}
		}

	})
	c.Start()

	return nil
}
