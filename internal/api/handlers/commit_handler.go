package handlers

import (
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/api/services"

	"github.com/gin-gonic/gin"
)

func GetCommitsByRepository(c *gin.Context) {
	repoName := c.Param("repoName")

	pageStr := c.Query("page")
	sizeStr := c.Query("page_size")

	commits, err := services.FetchCommitsByRepository(repoName, pageStr, sizeStr)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Unable to fetch commits"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": commits, "message": "Commits fetched"})
}

func GetTopNCommitAuthors(c *gin.Context) {
	n := c.Query("n")

	commits, err := services.FetchTopNCommitAuthors(n)
	if err != nil {
		c.JSON(http.StatusOK, gin.H{"error": "Unable to fetch top N commit authors"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": commits, "message": "Top N commit authors fetched"})
}
