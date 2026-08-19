package user

import (
	"errors"
	"net/http"

	"omniforge-api/internal/requestcontext"
	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	userService *Service
}

func NewHandler(userService *Service) *Handler {
	return &Handler{
		userService: userService,
	}
}

func (h *Handler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(requestcontext.UserID)

	if !exists {
		response.Error(
			c,
			http.StatusUnauthorized,
			MessageAuthenticationRequired,
		)
		return
	}

	userID, ok := userIDValue.(uint)
	if !ok {
		response.Error(
			c,
			http.StatusUnauthorized,
			MessageAuthenticationRequired,
		)
		return
	}

	result, err := h.userService.GetProfile(userID)
	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				MessageUserNotFound,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			MessageProfileFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageProfileSuccess,
		result,
	)
}