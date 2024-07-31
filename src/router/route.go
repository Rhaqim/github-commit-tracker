package router

import (
	"savannahtech/src/core"
	"savannahtech/src/model"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	// Set up the Gin Gonic server
	r := gin.Default()

	// Define your API routes here
	r.GET("/api/repositories/:owner/:repo", func(c *gin.Context) {
		owner := c.Param("owner")
		repo := c.Param("repo")
		if err := core.RepositoryData(owner, repo); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, gin.H{"status": "Repository data fetched"})
	})

	// r.GET("/api/commits/:owner/:repo", func(c *gin.Context) {
	// 	owner := c.Param("owner")
	// 	repo := c.Param("repo")
	// 	if err := core.CommitData(owner, repo); err != nil {
	// 		c.JSON(500, gin.H{"error": err.Error()})
	// 		return
	// 	}
	// 	c.JSON(200, gin.H{"status": "Repository data fetched"})
	// })

	r.GET("api/commits", func(c *gin.Context) {
		var commit model.CommitStore

		if err := commit.GetCommits(); err != nil {
			c.JSON(500, gin.H{"error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"commits": commit})

	})

	return r
}
