package main

import (
	"net/http"

	"omniforge-api/internal/config"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "OmniForge API is running",
			"environment": cfg.AppEnv,
		})
	})

	router.Run(":" + cfg.AppPort)
}