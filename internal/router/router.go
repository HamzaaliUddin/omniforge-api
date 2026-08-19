package router

import (
	"net/http"

	"omniforge-api/internal/app"
	"omniforge-api/internal/auth"

	"github.com/gin-gonic/gin"
)

func New(dependencies *app.Dependencies) *gin.Engine {
	router := gin.Default()

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status": "success",
			"message": "OmniForge API is running",
		})
	})

	api := router.Group("/api/v1")

	auth.RegisterRoutes(
		api,
		dependencies.AuthHandler,
	)

	return router
}