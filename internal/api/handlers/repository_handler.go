package handlers

import (
	"net/http"

	"github.com/Rhaqim/savannahtech/internal/api/services"

	"github.com/gin-gonic/gin"
)

func ProcessRepository(c *gin.Context) {
	owner := c.Param("owner")
	repo := c.Param("repo")
	err := services.ProcessRepository(owner, repo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to process repository"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Repository processed"})
}

func GetRepositoryByOwnerRepo(c *gin.Context) {
	ownerRepo := c.Param("owner_repo")
	repo, err := services.FetchRepositoryByOwnerRepo(ownerRepo)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to fetch repository"})
		return
	}
	c.JSON(http.StatusOK, repo)
}
