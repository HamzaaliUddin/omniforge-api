package admin

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	router.GET("/users", handler.GetUsers)
}