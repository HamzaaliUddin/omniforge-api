package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const permissionCacheTTL = 5 * time.Minute

func permissionKey(userID string) string {
	return fmt.Sprintf("rbac:permissions:%s", userID)
}

func GetPermissions(
	ctx context.Context,
	client *redis.Client,
	userID string,
) ([]string, bool, error) {
	data, err := client.Get(
		ctx,
		permissionKey(userID),
	).Result()

	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf("failed to get permissions cache: %w", err)
	}

	var permissions []string

	if err := json.Unmarshal([]byte(data), &permissions); err != nil {
		return nil, false, fmt.Errorf("failed to decode permissions cache: %w", err)
	}

	return permissions, true, nil
}

func SetPermissions(
	ctx context.Context,
	client *redis.Client,
	userID string,
	permissions []string,
) error {
	data, err := json.Marshal(permissions)
	if err != nil {
		return fmt.Errorf("failed to encode permissions cache: %w", err)
	}

	if err := client.Set(
		ctx,
		permissionKey(userID),
		data,
		permissionCacheTTL,
	).Err(); err != nil {
		return fmt.Errorf("failed to set permissions cache: %w", err)
	}

	return nil
}

func DeletePermissions(
	ctx context.Context,
	client *redis.Client,
	userID string,
) error {
	if err := client.Del(
		ctx,
		permissionKey(userID),
	).Err(); err != nil {
		return fmt.Errorf("failed to delete permissions cache: %w", err)
	}

	return nil
}