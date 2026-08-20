package ai

import (
	"net/http"

	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(
	service *Service,
) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) GenerateText(c *gin.Context) {
	var input GenerateTextRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	result, err := h.service.GenerateText(
		c.Request.Context(),
		input,
	)

	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			MessageGenerateFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageGenerateSuccess,
		result,
	)
}