package router

import (
	"omniforge-api/internal/admin"
	"omniforge-api/internal/app"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/middleware"
	"omniforge-api/internal/user"

	"github.com/gin-gonic/gin"
)

func New(dependencies *app.Dependencies) *gin.Engine {
	router := gin.Default()

	api := router.Group("/api/v1")

	auth.RegisterRoutes(
		api,
		dependencies.AuthHandler,
	)

	protectedRoutes := api.Group("")
	protectedRoutes.Use(
		middleware.Auth(dependencies.JWTSecret),
	)

	user.RegisterRoutes(
		protectedRoutes,
		dependencies.UserHandler,
	)

	adminRoutes := api.Group("/admin")
	adminRoutes.Use(
		middleware.Auth(dependencies.JWTSecret),
		middleware.RequireRole("admin"),
	)

	admin.RegisterRoutes(
		adminRoutes,
		dependencies.AdminHandler,
	)

	return router
}