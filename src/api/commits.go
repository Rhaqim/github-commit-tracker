package api

import (
	"savannahtech/src/core"
	"savannahtech/src/model"

	"github.com/gin-gonic/gin"
)

func GetCommit(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	if err := core.CommitData(owner, repo); err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}
	c.JSON(200, gin.H{"status": "Repository data fetched"})
}

func GetCommits(c *gin.Context) {
	var commit model.CommitStore

	commits, err := commit.GetCommits()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"commits": commits})

}
