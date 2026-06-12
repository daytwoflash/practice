package middlewares

import (
	"net/http"
	"project/logger"
	"project/utils"

	"github.com/gin-gonic/gin"
)

func Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				logger.Error("panic recovered",
					"panic", r,
					"method", c.Request.Method,
					"path", c.Request.URL.Path,
				)
				utils.Error(c, http.StatusInternalServerError, "Internal server error")
				c.Abort()
			}
		}()

		c.Next()
	}
}
