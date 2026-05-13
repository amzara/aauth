package redis

import (
	"context"
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *redis.Client
}

func NewRedisService(ctx context.Context, addr string, pw string) (*RedisService, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("redis ping failed: $w", err)
	}
	return &RedisService{
		Client: rdb,
	}, nil

}
