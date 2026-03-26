package handler

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func logrusMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		query := c.Request.URL.RawQuery
		c.Next()

		entry := logrus.WithFields(logrus.Fields{
			"status":     c.Writer.Status(),
			"method":     c.Request.Method,
			"path":       path,
			"query":      query,
			"latency_ms": time.Since(start).Milliseconds(),
			"client_ip":  c.ClientIP(),
		})
		if len(c.Errors) > 0 {
			entry.Error(c.Errors.String())
		} else {
			entry.Info("request completed")
		}
	}
}

func recoveryMiddleware() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		logrus.WithField("error", recovered).Error("panic recovered")
		c.AbortWithStatus(500)
	})
}
