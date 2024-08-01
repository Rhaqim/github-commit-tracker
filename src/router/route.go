package router

import (
	"savannahtech/src/api"

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
		commitGroup.GET("/get", api.GetCommits)
		commitGroup.GET("/top", api.GetTopCommitAuthors)
	}

	return r
}
