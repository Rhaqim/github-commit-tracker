package router

import (
	"github.com/Rhaqim/savannahtech/internal/api"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {

	r := gin.Default()

	repoGroup := r.Group("/repositories")
	{
		repoGroup.GET("/get/:owner/:repo", api.GetRepo)
		repoGroup.GET("/:owner/:repo", api.ProcessRepo)
	}

	commitGroup := r.Group("/commits")
	{
		commitGroup.GET("/top-authors", api.GetTopCommitAuthors)
		commitGroup.GET("/:repoName", api.GetCommitsByRepo)
	}

	return r
}
