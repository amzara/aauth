package session

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

const sessionTTL = 2 * time.Hour
const sessionPrefix = "session:"

var ErrSessionNotFound = errors.New("session does not exist")

type Store struct {
	rdb *redis.Client
}

func NewStore(rdb *redis.Client) *Store {
	return &Store{
		rdb: rdb,
	}
}

func generateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func (s *Store) Create(ctx context.Context, userID string, data map[string]string) (string, error) {
	token, err := generateToken()
	if err != nil {
		return "", err
	}
	key := sessionPrefix + token
	fields := make([]interface{}, 0, len(data)*2+2)
	fields = append(fields, "userID", userID)
	for k, v := range data {
		fields = append(fields, k, v)
	}

	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, key, fields...)
	pipe.Expire(ctx, key, sessionTTL)
	_, err = pipe.Exec(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to create session %w", err)
	}

	return token, nil

}

func (s *Store) Get(ctx context.Context, token string) (map[string]string, error) {
	result, err := s.rdb.HGetAll(ctx, sessionPrefix+token).Result()
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, ErrSessionNotFound
	}
	return result, nil
}

func (s *Store) Refresh(ctx context.Context, token string) error {
	return s.rdb.Expire(ctx, sessionPrefix+token, sessionTTL).Err()
}

func (s *Store) Set(ctx context.Context, token, field, value string) error {
	pipe := s.rdb.Pipeline()
	pipe.HSet(ctx, sessionPrefix+token, field, value)
	_, err := pipe.Exec(ctx)
	return err

}

func (s *Store) Destroy(ctx context.Context, token string) error {
	return s.rdb.Del(ctx, sessionPrefix+token).Err()
}

func (s *Store) Check(ctx context.Context, token string) (bool, error) {
	n, err := s.rdb.Exists(ctx, sessionPrefix+token).Result()
	return n > 0, err
}
