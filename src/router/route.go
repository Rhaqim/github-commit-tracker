package router

import (
	"savannahtech/src/api"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	// Set up the Gin Gonic server
	r := gin.Default()

	// Define your API routes here
	r.GET("/api/repositories/:owner/:repo", api.GetRepo)

	r.GET("/api/commits/:owner/:repo", api.GetCommit)
	r.GET("api/commits", api.GetCommits)

	return r
}
