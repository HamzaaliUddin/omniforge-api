package cache

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type RefreshSession struct {
	UserID uint   `json:"user_id"`
	Role   string `json:"role"`
}

func refreshTokenKey(token string) string {
	hash := sha256.Sum256([]byte(token))

	return fmt.Sprintf(
		"auth:refresh:%x",
		hash,
	)
}

func StoreRefreshToken(
	ctx context.Context,
	client *redis.Client,
	token string,
	session RefreshSession,
	ttl time.Duration,
) error {
	data, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("failed to encode refresh session: %w", err)
	}

	if err := client.Set(
		ctx,
		refreshTokenKey(token),
		data,
		ttl,
	).Err(); err != nil {
		return fmt.Errorf("failed to store refresh token: %w", err)
	}

	return nil
}

func GetRefreshSession(
	ctx context.Context,
	client *redis.Client,
	token string,
) (*RefreshSession, bool, error) {
	data, err := client.Get(
		ctx,
		refreshTokenKey(token),
	).Bytes()

	if err == redis.Nil {
		return nil, false, nil
	}

	if err != nil {
		return nil, false, fmt.Errorf(
			"failed to get refresh token: %w",
			err,
		)
	}

	var session RefreshSession

	if err := json.Unmarshal(data, &session); err != nil {
		return nil, false, fmt.Errorf(
			"failed to decode refresh session: %w",
			err,
		)
	}

	return &session, true, nil
}

func DeleteRefreshToken(
	ctx context.Context,
	client *redis.Client,
	token string,
) error {
	if err := client.Del(
		ctx,
		refreshTokenKey(token),
	).Err(); err != nil {
		return fmt.Errorf(
			"failed to delete refresh token: %w",
			err,
		)
	}

	return nil
}