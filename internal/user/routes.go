package user

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	userRoutes := router.Group("/users")

	userRoutes.GET("/me", handler.Me)
	userRoutes.PATCH("/me", handler.UpdateMe)
	userRoutes.DELETE("/me", handler.DeleteMe)
}