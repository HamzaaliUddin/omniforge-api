package router

import (
	"time"

	"omniforge-api/internal/admin"
	"omniforge-api/internal/ai"
	"omniforge-api/internal/app"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/middleware"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

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

	auth.RegisterPublicRoutes(
		api,
		deps.AuthHandler,
	)

	protected := api.Group("")

	protected.Use(
		middleware.Auth(
			deps.JWTSecret,
			deps.Redis,
		),
	)

	auth.RegisterProtectedRoutes(
		protected,
		deps.AuthHandler,
	)

	user.RegisterRoutes(
		protected,
		deps.UserHandler,
	)

	adminRoutes := protected.Group("")

	adminRoutes.Use(
		middleware.RequireRole(
			role.NameAdmin,
		),
	)

	admin.RegisterRoutes(
		adminRoutes,
		deps.AdminHandler,
	)

	ai.RegisterRoutes(
	protected,
	deps.AIHandler,
)

	return router
}