package main

import (
	"log"

	"omniforge-api/internal/app"
	"omniforge-api/internal/cache"
	"omniforge-api/internal/config"
	"omniforge-api/internal/database"
	"omniforge-api/internal/router"
	"omniforge-api/internal/seed"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg.DatabaseURL)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}
	log.Println("Database connected successfully")

	if err := seed.Run(db); err != nil {
		log.Fatal("Seeding failed: ", err)
	}

	log.Println("Database seeded successfully")

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	redisClient, err := cache.ConnectRedis(cfg.RedisAddr)
	
	if err != nil {
		log.Fatal(err)
	}

		defer redisClient.Close()

	log.Println("Redis connected successfully")
	dependencies := app.NewDependencies(db, cfg,redisClient)

	appRouter := router.New(dependencies)

	if err := appRouter.Run(":" + cfg.AppPort); err != nil {
		log.Fatal(err)
	}
}