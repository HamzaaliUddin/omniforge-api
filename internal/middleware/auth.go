package middleware

import (
	"net/http"
	"strings"

	"omniforge-api/internal/auth"
	"omniforge-api/internal/cache"
	"omniforge-api/internal/requestcontext"
	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

const (
	MessageTokenRevoked             = "token has been revoked"
	MessageAuthenticationUnavailable = "authentication service unavailable"
)

func Auth(
	jwtSecret string,
	redisClient *redis.Client,
) gin.HandlerFunc {
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

		tokenString := parts[1]

		claims, err := auth.ValidateAccessToken(
			tokenString,
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

		blacklisted, err := cache.IsTokenBlacklisted(
			c.Request.Context(),
			redisClient,
			tokenString,
		)

		if err != nil {
			response.Error(
				c,
				http.StatusServiceUnavailable,
				MessageAuthenticationUnavailable,
			)
			c.Abort()
			return
		}

		if blacklisted {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageTokenRevoked,
			)
			c.Abort()
			return
		}

		c.Set(requestcontext.UserID, claims.UserID)
		c.Set(requestcontext.Role, claims.Role)

		c.Next()
	}
}