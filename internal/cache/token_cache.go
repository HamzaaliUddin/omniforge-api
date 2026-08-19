package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

func tokenBlacklistKey(token string) string {
	return fmt.Sprintf("auth:blacklist:%s", token)
}

func BlacklistToken(
	ctx context.Context,
	client *redis.Client,
	token string,
	ttl time.Duration,
) error {
	if err := client.Set(
		ctx,
		tokenBlacklistKey(token),
		"1",
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("failed to blacklist token: %w", err)
	}

	return nil
}

func IsTokenBlacklisted(
	ctx context.Context,
	client *redis.Client,
	token string,
) (bool, error) {
	result, err := client.Exists(
		ctx,
		tokenBlacklistKey(token),
	).Result()

	if err != nil {
		return false, fmt.Errorf("failed to check token blacklist: %w", err)
	}

	return result > 0, nil
}