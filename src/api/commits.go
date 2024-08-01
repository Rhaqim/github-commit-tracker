package api

import (
	"savannahtech/src/model"

	"github.com/gin-gonic/gin"
)

func GetCommits(c *gin.Context) {
	var commit model.CommitStore

	commits, err := commit.GetCommits()
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"commits": commits})

}
