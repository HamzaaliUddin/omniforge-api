package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	adminRoutes := router.Group("/admin")

	adminRoutes.GET("/users", handler.GetUsers)
}