package router

import (
	"omniforge-api/internal/admin"
	"omniforge-api/internal/app"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/middleware"
	"omniforge-api/internal/user"
	"time"

	"github.com/gin-gonic/gin"
)

func New(deps *app.Dependencies) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")
	api.Use(
		middleware.RateLimit(
			deps.Redis,
			100,
			time.Minute,
		),
	)

	auth.RegisterRoutes(
		api,
		deps.AuthHandler,
	)

	protectedRoutes := api.Group("")
	protectedRoutes.Use(
		middleware.Auth(
		deps.JWTSecret,
		deps.Redis,
		),
	)

	user.RegisterRoutes(
		protectedRoutes,
		deps.UserHandler,
	)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.Auth(
		deps.JWTSecret,
		deps.Redis,
		),
		middleware.RequireRole("admin"),
	)

	admin.RegisterRoutes(
		adminRoutes,
		deps.AdminHandler,
	)

	return router
}