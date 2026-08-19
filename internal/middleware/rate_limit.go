package middleware

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"omniforge-api/internal/response"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)


var rateLimitScript = redis.NewScript(`
	local current = redis.call("INCR", KEYS[1])

	if current == 1 then
		redis.call("EXPIRE", KEYS[1], ARGV[1])
	end

	return current
`)

func RateLimit(
	redisClient *redis.Client,
	limit int64,
	window time.Duration,
) gin.HandlerFunc {
	return func(c *gin.Context) {
		path := c.FullPath()

		if path == "" {
			path = c.Request.URL.Path
		}

		key := fmt.Sprintf(
			"rate_limit:%s:%s:%s",
			c.ClientIP(),
			c.Request.Method,
			path,
		)

		count, err := rateLimitScript.Run(
			c.Request.Context(),
			redisClient,
			[]string{key},
			int64(window.Seconds()),
		).Int64()

		if err != nil {
			log.Printf("rate limiter redis error: %v", err)
			c.Next()
			return
		}

		if count > limit {
			response.Error(
				c,
				http.StatusTooManyRequests,
				MessageTooManyRequests,
			)

			c.Abort()
			return
		}

		c.Next()
	}
}