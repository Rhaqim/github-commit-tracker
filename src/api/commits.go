package api

import (
	"savannahtech/src/model"
	"strconv"

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

func GetTopCommitAuthors(c *gin.Context) {
	var commit model.CommitStore

	topN := c.Query("topN")

	topNInt, err := strconv.Atoi(topN)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	commitCounts, err := commit.GetTopCommitAuthors(topNInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	c.JSON(200, gin.H{"commitCounts": commitCounts})
}
