package middleware

import (
	"context"
	"fmt"
	"net/http"

	"omniforge-api/internal/cache"
	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type PermissionLoader func(
	ctx context.Context,
	userID string,
) ([]string, error)

func RequirePermission(
	redisClient *redis.Client,
	loadPermissions PermissionLoader,
	requiredPermission string,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		userIDValue, exists := c.Get("userID")
		if !exists {
			response.Error(
				c,
				http.StatusUnauthorized,
				MessageUnauthorized,
			)
			c.Abort()
			return
		}

		userID := fmt.Sprint(userIDValue)

		permissions, found, err := cache.GetPermissions(
			c.Request.Context(),
			redisClient,
			userID,
		)

		if err != nil {
			found = false
		}

		if !found {
			permissions, err = loadPermissions(
				c.Request.Context(),
				userID,
			)

			if err != nil {
				response.Error(
					c,
					http.StatusInternalServerError,
					MessageServerError,
				)
				c.Abort()
				return
			}

			_ = cache.SetPermissions(
				c.Request.Context(),
				redisClient,
				userID,
				permissions,
			)
		}

		if !hasPermission(permissions, requiredPermission) {
			response.Error(
				c,
				http.StatusForbidden,
				MessageForbidden,
			)
			c.Abort()
			return
		}

		c.Next()
	}
}

func hasPermission(
	permissions []string,
	requiredPermission string,
) bool {
	for _, permission := range permissions {
		if permission == requiredPermission {
			return true
		}
	}

	return false
}