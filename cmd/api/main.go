package main

import (
	"log"
	"net/http"

	"omniforge-api/internal/config"
	"omniforge-api/internal/database"

	"github.com/gin-gonic/gin"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	router := gin.Default()

	router.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message":     "OmniForge API is running",
			"environment": cfg.AppEnv,
		})
	})

	if err := router.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}