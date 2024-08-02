package api

import (
	"savannahtech/src/core"
	"savannahtech/src/model"
	"strings"

	"github.com/gin-gonic/gin"
)

func ProcessRepo(c *gin.Context) {
	owner := strings.ToLower(c.Param("owner"))
	repo := strings.ToLower(c.Param("repo"))
	fromDate := c.Query("fromDate")

	if err := core.ProcessRepositoryData(owner, repo, fromDate); err != nil {
		c.JSON(400, gin.H{"developer_error": err.Error(), "message": "failed to process repository data"})
		return
	}
	c.JSON(200, gin.H{"status": "Processed repository data"})
}

func GetRepo(c *gin.Context) {
	var repoStore model.RepositoryStore

	owner := strings.ToLower(c.Param("owner"))
	repo := strings.ToLower(c.Param("repo"))

	err := repoStore.GetRepositoryByOwnerRepo(owner + "/" + repo)
	if err != nil {
		c.JSON(500, gin.H{"developer_error": err.Error(), "message": "failed to retrieve repository"})
		return
	}

	c.JSON(200, gin.H{"repo": repoStore, "message": "repository retrieved successfully"})
}
