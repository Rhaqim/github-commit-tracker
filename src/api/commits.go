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
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to retrieve commits"})
		return
	}

	c.JSON(200, gin.H{"commits": commits, "message": "commits retrieved successfully"})

}

func GetTopCommitAuthors(c *gin.Context) {
	var commit model.CommitStore

	topN := c.Query("n")

	topNInt, err := strconv.Atoi(topN)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error(), "message": "invalid topN value"})
		return
	}

	commitCounts, err := commit.GetTopCommitAuthors(topNInt)
	if err != nil {
		c.JSON(500, gin.H{"error": err.Error(), "message": "failed to retrieve top commit authors"})
		return
	}

	c.JSON(200, gin.H{"commitCounts": commitCounts, "message": "top commit authors retrieved successfully"})
}
