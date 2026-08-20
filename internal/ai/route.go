package ai

import "github.com/gin-gonic/gin"

func RegisterRoutes(
	router *gin.RouterGroup,
	handler *Handler,
) {
	aiRoutes := router.Group("/ai")

	aiRoutes.POST(
		"/text/generate",
		handler.GenerateText,
	)
}