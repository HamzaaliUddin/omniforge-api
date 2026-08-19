package app

import (
	"omniforge-api/internal/admin"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/config"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	DB *gorm.DB
	Redis *redis.Client
	AuthHandler  *auth.Handler
	UserHandler  *user.Handler
	AdminHandler *admin.Handler
	JWTSecret    string
}

func NewDependencies(db *gorm.DB,cfg *config.Config,redisClient *redis.Client) *Dependencies {
	userRepository := user.NewRepository(db)
	roleRepository := role.NewRepository(db)

	authService := auth.NewService(
		userRepository,
		roleRepository,
		redisClient,
		cfg.JWTSecret,
	)

	authHandler := auth.NewHandler(authService)

	userService := user.NewService(userRepository)
	userHandler := user.NewHandler(userService)

	adminService := admin.NewService(userRepository)
	adminHandler := admin.NewHandler(adminService)

return &Dependencies{
	DB:    db,
	Redis: redisClient,
	AuthHandler:  authHandler,
	UserHandler:  userHandler,
	AdminHandler: adminHandler,
	JWTSecret:    cfg.JWTSecret,
}
}