/* 后端内部错误 */

package utils

import (
	"errors"
	"fmt"
	"project/logger"

	"github.com/gin-gonic/gin"
)

type AppError struct {
	Code    int
	Message string
	Err     error
}

// 实现Error()接口，*AppError自动“继承”标准库中的error
func (e *AppError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message

}

func NewAppError(code int, message string) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
	}
}

// 业务错误AppError处理
func HandleError(c *gin.Context, err error) {
	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Code >= 500 {
			logger.Error("server error",
				"code", appErr.Code,
				"message", appErr.Message,
				"err", appErr.Err,
				"path", c.Request.URL.Path)
		} else {
			logger.Warn("client error",
				"code", appErr.Code,
				"message", appErr.Message,
				"path", c.Request.URL.Path)
		}
		Error(c, appErr.Code, appErr.Message)
		return
	}

	// 未知错误，记录日志但不暴露给客户端
	logger.Error("unxepected error", "err", err, "path", c.Request.URL.Path)
	Error(c, 500, "Internal server error")
}
