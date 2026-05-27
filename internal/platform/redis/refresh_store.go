package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	goredis "github.com/redis/go-redis/v9"

	authapp "github.com/mohit838/olario-platform-backend/internal/application/auth"
)

// RefreshStore stores refresh token sessions in Redis.
// Refresh tokens are rotated by deleting the old token hash and storing a new
// hash, so replaying an old refresh token fails.
type RefreshStore struct {
	client *goredis.Client
}

func NewRefreshStore(client *goredis.Client) *RefreshStore {
	return &RefreshStore{client: client}
}

func (s *RefreshStore) StoreRefresh(ctx context.Context, tokenHash string, session authapp.RefreshSession, ttl time.Duration) error {
	payload, err := json.Marshal(session)
	if err != nil {
		return fmt.Errorf("marshal refresh session: %w", err)
	}
	return s.client.Set(ctx, refreshKey(tokenHash), payload, ttl).Err()
}

func (s *RefreshStore) RotateRefresh(ctx context.Context, oldHash, newHash string, ttl time.Duration) (authapp.RefreshSession, error) {
	oldKey := refreshKey(oldHash)
	raw, err := s.client.GetDel(ctx, oldKey).Result()
	if err != nil {
		return authapp.RefreshSession{}, fmt.Errorf("read old refresh session: %w", err)
	}

	var session authapp.RefreshSession
	if err := json.Unmarshal([]byte(raw), &session); err != nil {
		return authapp.RefreshSession{}, fmt.Errorf("parse refresh session: %w", err)
	}

	if err := s.StoreRefresh(ctx, newHash, session, ttl); err != nil {
		return authapp.RefreshSession{}, err
	}
	return session, nil
}

func (s *RefreshStore) DeleteRefresh(ctx context.Context, tokenHash string) error {
	return s.client.Del(ctx, refreshKey(tokenHash)).Err()
}

func refreshKey(tokenHash string) string {
	return "auth:refresh:" + tokenHash
}
