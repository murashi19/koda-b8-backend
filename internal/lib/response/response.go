package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Result  any    `json:"result,omitempty"`
}

func Success(c *gin.Context, status int, message string, result interface{}) {
	c.JSON(status, Response{
		Success: true,
		Message: message,
		Result:  result,
	})
}

func Error(c *gin.Context, status int, message string) {
	c.JSON(status, Response{
		Success: false,
		Message: message,
	})
}

func InternalServerError(c *gin.Context) {
	Error(c, http.StatusInternalServerError, "Internal server error")
}

func Unauthorized(c *gin.Context) {
	Error(c, http.StatusUnauthorized, "Unauthorized")
}

func Forbidden(c *gin.Context) {
	Error(c, http.StatusForbidden, "Forbidden")
}

func NotFound(c *gin.Context) {
	Error(c, http.StatusNotFound, "Resource not found")
}

func BadRequest(c *gin.Context, message string) {
	Error(c, http.StatusBadRequest, message)
}
