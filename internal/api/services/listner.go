package services

import (
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func ProcessFunc(event entities.Event) error {
	owner, repo, startDate := event.Owner, event.Repo, event.StartDate

	// Process repository event
	if event.Type == entities.RepoEvent {
		go func() {
			if err := ProcessRepository(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process repository:", err)
			}
		}()
	}

	// Process commit event
	if event.Type == entities.CommitEvent {
		go func() {
			if err := ProcessCommitData(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process commit data:", err)
			}
		}()
	}

	// Process periodic fetch event
	if event.Type == entities.PeriodEvent {
		go func() {
			if err := PeriodicFetch(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process periodic fetch:", err)
			}
		}()
	}

	return nil
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
