package middleware

import (
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// Logs server requests
func Logger(logger *zap.SugaredLogger) gin.HandlerFunc {
	return func(c *gin.Context) {
		t := time.Now()

		// Request
		c.Next()

		// Counting request duration
		latency := time.Since(t)
		miliSeconds := fmt.Sprintf("%d ms", latency.Milliseconds())

		// Status code and body size
		status := c.Writer.Status()
		size := c.Writer.Size()

		// Log body
		logParams := []interface{}{
			"method", c.Request.Method,
			"uri", c.Request.URL,
			"status", status,
			"duration", miliSeconds,
			"size", size,
		}

		// Write to logs
		if status >= 500 {
			logger.Errorln(logParams...)
		} else {
			logger.Infoln(logParams...)
		}
	}
}
