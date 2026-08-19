package middleware

import (
	"net/http"

	"omniforge-api/internal/response"
	"omniforge-api/internal/requestcontext"

	"github.com/gin-gonic/gin"
)

func RequireRole(allowedRoles ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		roleValue, exists := c.Get(requestcontext.Role)

		if !exists {
			response.Error(
				c,
				http.StatusForbidden,
				MessageAccessDenied,
			)
			c.Abort()
			return
		}

		role, ok := roleValue.(string)
		if !ok {
			response.Error(
				c,
				http.StatusForbidden,
				MessageAccessDenied,
			)
			c.Abort()
			return
		}

		for _, allowedRole := range allowedRoles {
			if role == allowedRole {
				c.Next()
				return
			}
		}

		response.Error(
			c,
			http.StatusForbidden,
			MessageAccessDenied,
		)
		c.Abort()
	}
}