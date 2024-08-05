package services

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func ProcessFunc(event entities.Event) error {
	var owner, repo, startDate string = event.Owner, event.Repo, event.StartDate

	if event.Type == entities.RepoEvent {
		go func() {
			if err := ProcessRepository(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process repository:", err)
			}
		}()
	}
	if event.Type == entities.CommitEvent {
		go func() {
			if err := ProcessCommitData(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process commit data:", err)
			}
		}()
	}

	if event.Type == entities.PeriodEvnt {
		go func() {
			if err := PeriodicFetch(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process periodic fetch:", err)
			}

		}()
	}

	// if event.Type == entities.RepoEvent {
	// 	return ProcessRepository(owner, repo, startDate)
	// }
	// if event.Type == entities.CommitEvent {
	// 	return ProcessCommitData(owner, repo, startDate)
	// }
	// if event.Type == entities.PeriodEvnt {
	// 	return PeriodicFetch(owner, repo, startDate)
	// }

	return nil
}
