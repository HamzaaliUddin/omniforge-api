package middleware

import (
	"net/http"
	"strings"

	"omniforge-api/internal/auth"
	"omniforge-api/internal/requestcontext"
	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
)

func Auth(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorization := c.GetHeader("Authorization")

		if authorization == "" {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageAuthenticationRequired,
			)
			c.Abort()
			return
		}

		parts := strings.SplitN(authorization, " ", 2)

		if len(parts) != 2 ||
			!strings.EqualFold(parts[0], "Bearer") ||
			parts[1] == "" {

			response.Error(
				c,
				http.StatusUnauthorized,
				MessageInvalidToken,
			)
			c.Abort()
			return
		}

		claims, err := auth.ValidateAccessToken(
			parts[1],
			jwtSecret,
		)

		if err != nil {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageInvalidToken,
			)
			c.Abort()
			return
		}

		c.Set(requestcontext.UserID, claims.UserID)
		c.Set(requestcontext.Role, claims.Role)

		c.Next()
	}
}