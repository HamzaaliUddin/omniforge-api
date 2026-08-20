package user

import (
	"errors"
	"net/http"

	"omniforge-api/internal/requestcontext"
	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{
		service: service,
	}
}

func (h *Handler) Me(c *gin.Context) {
	userIDValue, exists := c.Get(
		requestcontext.UserID,
	)

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

	result, err := h.service.GetMe(
		c.Request.Context(),
		userID,
	)

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

func (h *Handler) UpdateMe(c *gin.Context) {
	userIDValue, exists := c.Get(
		requestcontext.UserID,
	)

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

	var input UpdateMeRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	result, err := h.service.UpdateMe(
		c.Request.Context(),
		userID,
		input,
	)

	if err != nil {
		if errors.Is(err, ErrUserNotFound) {
			response.Error(
				c,
				http.StatusNotFound,
				MessageUserNotFound,
			)
			return
		}

		if errors.Is(err, ErrEmailAlreadyExists) {
			response.Error(
				c,
				http.StatusConflict,
				MessageEmailExists,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			MessageUpdateFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageUpdateSuccess,
		result,
	)
}

func (h *Handler) DeleteMe(c *gin.Context) {
	userIDValue, exists := c.Get(
		requestcontext.UserID,
	)

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

	err := h.service.DeleteMe(
		c.Request.Context(),
		userID,
	)

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
			MessageDeleteFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageDeleteSuccess,
		nil,
	)
}