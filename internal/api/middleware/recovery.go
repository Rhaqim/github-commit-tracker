package middleware

import (
	"net/http"

	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/gin-gonic/gin"
)

func RecoveryMiddleware() gin.HandlerFunc {
	return gin.RecoveryWithWriter(gin.DefaultErrorWriter, func(c *gin.Context, err interface{}) {
		logger.ErrorLogger.Printf("panic occurred: %s\n", err)
		c.JSON(http.StatusOK, gin.H{"error": "Internal Server Error"})
	})
}
