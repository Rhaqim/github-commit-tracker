package listener

import (
	"strings"

	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/utils"
	"github.com/Rhaqim/savannahtech/pkg/logger"
)

/*
ProcessFunc processes the event based on the event type.
*/
func ProcessFunc(event entities.Event) error {
	owner_, repo_, startDate_ := event.Owner, event.Repo, event.StartDate

	owner, repo := strings.ToLower(owner_), strings.ToLower(repo_)

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
			if err := services.PeriodicFetch(owner, repo); err != nil {
				logger.ErrorLogger.Println("Failed to process periodic fetch:", err)
			}
		}()
	}

	return nil
}
