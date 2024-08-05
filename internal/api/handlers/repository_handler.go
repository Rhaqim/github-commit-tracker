package handlers

import (
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/pkg/logger"

	"github.com/gin-gonic/gin"
)

func ProcessRepository(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	startDate := c.Query("start_date")

	// err := services.ProcessRepository(owner, repo, startDate)
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process repository"})
	// 	return
	// }

	// Send the event to the channel asynchronously
	go func() {
		err := services.ProcessRepository(owner, repo, startDate)
		if err != nil {
			logger.ErrorLogger.Println("Failed to process repository:", err)
		}
	}()

	// event := entities.Event{
	// 	Owner:     owner,
	// 	Repo:      repo,
	// 	StartDate: startDate,
	// 	Type:      entities.RepoEvent,
	// }

	// events.SendEvent(event)

	c.JSON(http.StatusOK, gin.H{"message": "Repository processed"})
}

func GetRepositoryByOwnerRepo(c *gin.Context) {
	ownerRepo := c.Param("owner_repo")
	repo, err := services.FetchRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch repository"})
		return
	}
	c.JSON(http.StatusOK, repo)
}
