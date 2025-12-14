package apperror

import (
	"github.com/gin-gonic/gin"
	"time"
)

type ErrorPayload struct {
	Code      string      `json:"code"`
	Message   string      `json:"message"`
	Details   interface{} `json:"details,omitempty"`
	Timestamp time.Time   `json:"timestamp"`
}

func JSONError(c *gin.Context, status int, code, message string, details interface{}) {
	c.AbortWithStatusJSON(status, ErrorPayload{Code: code, Message: message, Details: details, Timestamp: time.Now().UTC()})
}
