package app

import (
	"omniforge-api/internal/auth"
	"omniforge-api/internal/config"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler *auth.Handler
}

func NewDependencies(db *gorm.DB,cfg *config.Config,) *Dependencies {
	userRepository := user.NewRepository(db)
	roleRepository := role.NewRepository(db)

	authService := auth.NewService(
	userRepository,
	roleRepository,
	cfg.JWTSecret,
)

	authHandler := auth.NewHandler(authService)

	return &Dependencies{
		AuthHandler: authHandler,
	}
}