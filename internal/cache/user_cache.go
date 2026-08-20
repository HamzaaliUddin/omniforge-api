package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const userCacheTTL = 5 * time.Minute

type UserCache struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Role  string `json:"role"`
}

func userCacheKey(userID uint) string {
	return fmt.Sprintf("user:%d", userID)
}

func GetUser(
	ctx context.Context,
	client *redis.Client,
	userID uint,
) (*UserCache, bool, error) {
	data, err := client.Get(
		ctx,
		userCacheKey(userID),
	).Bytes()

	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf(
			"failed to get user cache: %w",
			err,
		)
	}

	var user UserCache

	if err := json.Unmarshal(data, &user); err != nil {
		return nil, false, fmt.Errorf(
			"failed to decode user cache: %w",
			err,
		)
	}

	return &user, true, nil
}

func SetUser(
	ctx context.Context,
	client *redis.Client,
	user UserCache,
) error {
	data, err := json.Marshal(user)
	if err != nil {
		return fmt.Errorf(
			"failed to encode user cache: %w",
			err,
		)
	}

	if err := client.Set(
		ctx,
		userCacheKey(user.ID),
		data,
		userCacheTTL,
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to set user cache: %w",
			err,
		)
	}

	return nil
}

func DeleteUser(
	ctx context.Context,
	client *redis.Client,
	userID uint,
) error {
	if err := client.Del(
		ctx,
		userCacheKey(userID),
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to delete user cache: %w",
			err,
		)
	}

	return nil
}