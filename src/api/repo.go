package api

import (
	"savannahtech/src/core"

	"github.com/gin-gonic/gin"
)

func GetRepo(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	if err := core.RepositoryData(owner, repo); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "Repository data fetched"})
}
