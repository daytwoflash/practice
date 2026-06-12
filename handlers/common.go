package handlers

import (
	"net/http"
	"project/utils"

	"github.com/gin-gonic/gin"
)

func getAuthenticatedUserID(c *gin.Context) (uint, bool) {
	userID, exists := c.Get("userID")
	if !exists {
		utils.Error(c, http.StatusUnauthorized, "Authentication required")
		return 0, false
	}

	id, ok := userID.(uint)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "Authentication required")
		return 0, false
	}

	return id, true

}

func parseValidationErrors(err error) map[string]string {
	errors := make(map[string]string)
	// 简化处理，实际应该解析 binding 错误
	errors["general"] = err.Error()
	return errors
}
