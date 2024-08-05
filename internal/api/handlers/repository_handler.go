package handlers

import (
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/api/services"
	"github.com/Rhaqim/savannahtech/internal/core/entities"
	"github.com/Rhaqim/savannahtech/internal/events"

	"github.com/gin-gonic/gin"
)

func ProcessRepository(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	startDate := c.Query("start_date")

	event := entities.Event{
		Owner:     owner,
		Repo:      repo,
		StartDate: startDate,
		Type:      entities.RepoEvent,
	}

	events.SendEvent(event)

	c.JSON(http.StatusOK, gin.H{"message": "Repository processing initiated"})
}

func GetRepositoryByOwnerRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	repo_, err := services.FetchRepositoryByOwnerRepo(owner, repo)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Unable to fetch repository"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": repo_, "message": "Repository fetched successfully"})
}
