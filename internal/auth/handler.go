package auth

import (
	"errors"
	"net/http"
	"strings"

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

func (h *Handler) Register(c *gin.Context) {
	var input RegisterRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	user, err := h.service.Register(input)
	if err != nil {
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
			MessageRegisterFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusCreated,
		MessageRegisterSuccess,
		gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email,
		},
	)
}

func (h *Handler) Login(c *gin.Context) {
	var input LoginRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	result, err := h.service.Login(
	c.Request.Context(),
	input,
	)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageInvalidCredentials,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			MessageLoginFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageLoginSuccess,
		result,
	)
}

func (h *Handler) Logout(c *gin.Context) {
	var input LogoutRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	authorization := c.GetHeader("Authorization")

	parts := strings.SplitN(authorization, " ", 2)

	if len(parts) != 2 ||
		!strings.EqualFold(parts[0], "Bearer") ||
		parts[1] == "" {

		response.Error(
			c,
			http.StatusUnauthorized,
			MessageTokenRequired,
		)
		return
	}

	err := h.service.Logout(
		c.Request.Context(),
		parts[1],
		input.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, ErrInvalidToken) {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageInvalidToken,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			MessageLogoutFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageLogoutSuccess,
		nil,
	)
}
func (h *Handler) Refresh(c *gin.Context) {
	var input RefreshRequest

	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			MessageInvalidRequest,
		)
		return
	}

	result, err := h.service.Refresh(
		c.Request.Context(),
		input.RefreshToken,
	)

	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageInvalidRefresh,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			MessageRefreshFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageRefreshSuccess,
		result,
	)
}