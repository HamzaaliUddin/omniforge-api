package app

import (
	"omniforge-api/internal/admin"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/config"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler  *auth.Handler
	UserHandler  *user.Handler
	AdminHandler *admin.Handler
	JWTSecret    string
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

userService := user.NewService(userRepository)
userHandler := user.NewHandler(userService)

adminService := admin.NewService(userRepository)
adminHandler := admin.NewHandler(adminService)

return &Dependencies{
	AuthHandler:  authHandler,
	UserHandler:  userHandler,
	AdminHandler: adminHandler,
	JWTSecret:    cfg.JWTSecret,
}
}