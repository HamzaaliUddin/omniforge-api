package admin

import (
	"net/http"

	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	adminService *Service
}

func NewHandler(adminService *Service) *Handler {
	return &Handler{
		adminService: adminService,
	}
}

func (h *Handler) GetUsers(c *gin.Context) {
	users, err := h.adminService.GetUsers()
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			MessageUsersFailed,
		)
		return
	}

	response.Success(
		c,
		http.StatusOK,
		MessageUsersSuccess,
		users,
	)
}