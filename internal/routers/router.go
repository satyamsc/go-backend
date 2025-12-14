package routers

import (
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
)

func New(_ *gorm.DB) *gin.Engine {
    r := gin.Default()
    r.GET("/health", func(c *gin.Context) { c.JSON(200, gin.H{"status": "ok"}) })
    return r
}
