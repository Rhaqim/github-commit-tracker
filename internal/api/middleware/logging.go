package middleware

import (
	"time"

	"github.com/Rhaqim/savannahtech/pkg/logger"
	"github.com/gin-gonic/gin"
)

func RequestLogger() gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		c.Next()
		endTime := time.Now()
		logger.ErrorLogger.Printf("[%s] %s %s %d %s\n", endTime.Format(time.RFC3339), c.Request.Method, c.Request.URL.Path, c.Writer.Status(), endTime.Sub(startTime))
	}
}
