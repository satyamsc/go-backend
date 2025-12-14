package utils

import "github.com/gin-gonic/gin"

func JSON(c *gin.Context, status int, v any) { c.JSON(status, v) }

