package router

import (
	"github.com/Rhaqim/savannahtech/internal/api/handlers"
	"github.com/Rhaqim/savannahtech/internal/api/middleware"
	"github.com/gin-gonic/gin"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	r.Use(middleware.RequestLogger())
	r.Use(middleware.RecoveryMiddleware())

	// // Define routes
	repoGroup := r.Group("/repositories")
	{
		repoGroup.GET("/process/:owner/:repo", handlers.ProcessRepository)
		repoGroup.GET("/get/:owner/:repo", handlers.GetRepositoryByOwnerRepo)
	}

	commitGroup := r.Group("/commits")
	{
		commitGroup.GET("/top-authors", handlers.GetTopNCommitAuthors)
		commitGroup.GET("/:repoName", handlers.GetCommitsByRepository)
	}

	r.GET("/health", handlers.HealthCheck)

	return r
}
