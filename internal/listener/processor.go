package listener

import (
	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

func ProcessFunc(event entities.Event) error {
	owner, repo, startDate_ := event.Owner, event.Repo, event.StartDate

	startDate := utils.ValidateDate(startDate_)

	// Process repository event
	if event.Type == entities.RepoEvent {
		go func() {
			if err := services.ProcessRepository(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process repository:", err)
			}
		}()
	}

	// Process commit event
	if event.Type == entities.CommitEvent {
		go func() {
			if err := services.ProcessCommitData(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process commit data:", err)
			}
		}()
	}

	// Process periodic fetch event
	if event.Type == entities.PeriodEvent {
		go func() {
			if err := services.PeriodicFetch(owner, repo, startDate); err != nil {
				logger.ErrorLogger.Println("Failed to process periodic fetch:", err)
			}
		}()
	}

	return nil
}
