package app

import (
	"omniforge-api/internal/admin"
	"omniforge-api/internal/ai"
	"omniforge-api/internal/auth"
	"omniforge-api/internal/role"
	"omniforge-api/internal/user"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Dependencies struct {
	AuthHandler  *auth.Handler
	UserHandler  *user.Handler
	AdminHandler *admin.Handler
	AIHandler    *ai.Handler

	Redis     *redis.Client
	JWTSecret string
}

func NewDependencies(
	db *gorm.DB,
	redisClient *redis.Client,
	jwtSecret string,
	openAIAPIKey string,
) *Dependencies {
	userRepository := user.NewRepository(db)
	roleRepository := role.NewRepository(db)

	userService := user.NewService(
		userRepository,
		redisClient,
	)

	authService := auth.NewService(
		userRepository,
		roleRepository,
		redisClient,
		jwtSecret,
	)

	adminService := admin.NewService(
		userRepository,
	)

	openAIProvider := ai.NewOpenAIProvider(
		openAIAPIKey,
	)

	aiService := ai.NewService(
		openAIProvider,
	)

	return &Dependencies{
		AuthHandler: auth.NewHandler(
			authService,
		),
		UserHandler: user.NewHandler(
			userService,
		),
		AdminHandler: admin.NewHandler(
			adminService,
		),
		AIHandler: ai.NewHandler(
			aiService,
		),

		Redis:     redisClient,
		JWTSecret: jwtSecret,
	}
}