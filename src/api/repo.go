package api

import (
	"savannahtech/src/config"
	"savannahtech/src/core"
	"savannahtech/src/event"
	"savannahtech/src/model"
	"savannahtech/src/types"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ProcessRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	if err := core.ProcessRepositoryData(owner, repo); err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to process repository data"})
		return
	}
	c.JSON(200, gin.H{"status": "Processed repository data"})
}

func GetRepo(c *gin.Context) {
	var repoStore model.RepositoryStore

	owner := c.Param("owner")
	repo := c.Param("repo")

	err := repoStore.GetRepositoriesByOwnerAndRepo(owner, repo)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to retrieve repository"})
		return
	}

	c.JSON(200, gin.H{"repo": repoStore, "message": "repository retrieved successfully"})
}

func ProcessSubRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")

	var newRepoEvent *event.EventQueue = event.NewEventQueue(config.NewRepo)
	if err := newRepoEvent.Publish(types.Event{
		ID:      uuid.New().String(),
		Message: "New repository event",
		Type:    types.NewRepo,
		Owner:   owner,
		Repo:    repo,
	}); err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to publish new repository event"})
		return
	}

	c.JSON(200, gin.H{"status": "Processed repository data"})
}
