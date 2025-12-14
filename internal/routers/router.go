package routers

import (
	"go-backend/internal/handlers"
	"go-backend/internal/middlewares"
	"go-backend/internal/repositories"
	"go-backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func New(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	r.Use(middlewares.CORS())
	r.Use(middlewares.GlobalRecovery())
	repo := repositories.NewDeviceRepository(db)
	svc := services.NewDeviceService(repo)
	h := handlers.NewDeviceHandler(svc)
	r.GET("/healthz", func(c *gin.Context) { c.Status(200) })
	r.GET("/openapi.yaml", func(c *gin.Context) { c.File("openapi.yaml") })
	r.GET("/docs", handlers.Docs)
	grp := r.Group("/devices")
	{
		grp.POST("", h.Create)
		grp.GET("", h.List)
		grp.GET("/:id", h.Get)
		grp.PUT("/:id", h.Update)
		grp.PATCH("/:id", h.Patch)
		grp.DELETE("/:id", h.Delete)
	}
	return r
}
