package core

import (
	"savannahtech/src/config"
	"savannahtech/src/log"
	"savannahtech/src/model"
	"savannahtech/src/types"
	"savannahtech/src/utils"
	"time"
)

func PeriodicFetch(owner, repo string) error {
	var commitStore model.CommitStore

	interval := 10 * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	ownerRepo := owner + "/" + repo
	baseURL := config.GithubRepoURL + ownerRepo + "/commits"

	// Start a goroutine to fetch commits periodically
	go func() {
		log.InfoLogger.Println("Fetching commits for " + ownerRepo + " ...")

		for range ticker.C {
			log.InfoLogger.Println("Running every " + interval.String())

			// Get the last commit date stored
			lastCommitDate := commitStore.GetLastCommitDate(ownerRepo)
			url := baseURL + "?since=" + lastCommitDate.String()

			commitsChan := make(chan []types.Commit)
			go func() {
				if err := utils.FetchCommits(url, commitsChan); err != nil {
					log.ErrorLogger.Printf("Failed to fetch commits: %v", err)
					close(commitsChan)
					return
				}
				close(commitsChan)
			}()

			for commit := range commitsChan {
				log.InfoLogger.Printf("Received %d commits for %s", len(commit), ownerRepo)
				if err := StoreCommit(commit, ownerRepo); err != nil {
					log.ErrorLogger.Printf("Error storing commits: %v", err)
				}
			}
		}
	}()

	return nil
}
