package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppPort     string
	AppEnv      string
	DatabaseURL string
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	return &Config{
		AppPort:     os.Getenv("APP_PORT"),
		AppEnv:      os.Getenv("APP_ENV"),
		DatabaseURL: os.Getenv("DATABASE_URL"),
	}
}