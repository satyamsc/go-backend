package middlewares

import (
	"fmt"
	"github.com/gin-gonic/gin"
	apperror "go-backend/pkg/error"
	"net/http"
)

func GlobalRecovery() gin.HandlerFunc {
	return gin.CustomRecovery(func(c *gin.Context, recovered interface{}) {
		apperror.JSONError(c, http.StatusInternalServerError, "internal_error", "unexpected error", map[string]string{"panic": fmt.Sprint(recovered)})
	})
}
