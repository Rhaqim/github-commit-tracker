package api

import (
	"strconv"
	"strings"

	"savannahtech/internal/model"
	"savannahtech/internal/utils"

	"github.com/gin-gonic/gin"
)

/*
GetTopCommitAuthors returns the top N commit authors in the database.

It retrieves the top N commit authors from the database and returns them as a JSON response.
*/
func GetTopCommitAuthors(c *gin.Context) {
	var commit model.CommitStore

	topN := c.Query("n")

	topNInt, err := strconv.Atoi(topN)
	if err != nil {
		c.JSON(400, gin.H{"developer_error": err.Error(), "message": "invalid topN value"})
		return
	}

	commitCounts, err := commit.GetTopCommitAuthors(topNInt)
	if err != nil {
		c.JSON(500, gin.H{"developer_error": err.Error(), "message": "failed to retrieve top commit authors"})
		return
	}

	c.JSON(200, gin.H{"commitCounts": commitCounts, "message": "top commit authors retrieved successfully"})
}

/*
GetCommitsByAuthor returns a list of commits for a repository by repository name.

It retrieves the commits for a repository from the database and returns them as a JSON response.
*/
func GetCommitsByRepo(c *gin.Context) {
	var commit model.CommitStore

	repoName := strings.ToLower(c.Param("repoName"))

	pageStr := c.Query("page")
	sizeStr := c.Query("page_size")

	size, offset := utils.Paigenation(pageStr, sizeStr)

	commits, err := commit.GetCommitsByAuthor(repoName, size, offset)
	if err != nil {
		c.JSON(500, gin.H{"developer_error": err.Error(), "message": "failed to retrieve commits for repository"})
		return
	}

	c.JSON(200, gin.H{"commits": commits, "message": "commits retrieved successfully"})
}
