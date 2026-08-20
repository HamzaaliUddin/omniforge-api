package auth

import "github.com/gin-gonic/gin"

func RegisterPublicRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	authRoutes := router.Group("/auth")

	authRoutes.POST("/register", handler.Register)
	authRoutes.POST("/login", handler.Login)
	authRoutes.POST("/refresh", handler.Refresh)
}

func RegisterProtectedRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	authRoutes := router.Group("/auth")

	authRoutes.POST("/logout", handler.Logout)
}