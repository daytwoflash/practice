package middlewares

import (
	"project/logger"
	"time"

	"github.com/gin-gonic/gin"
)

func Logger() gin.HandlerFunc {

	return func(c *gin.Context) {

		start := time.Now()
		path := c.Request.URL.Path
		method := c.Request.Method

		c.Next()

		latency := time.Since(start)
		status := c.Writer.Status()
		logger.Info("request completed",
			"method", method,
			"path", path,
			"status", status,
			"latency", latency,
			"client_ip", c.ClientIP(),
		)

	}

}
