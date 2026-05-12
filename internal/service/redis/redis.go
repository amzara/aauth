package redis

import (
	"context"
	"log"

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
	defer rdb.Close()
	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("Failed to connect to redis")
	}
	return &RedisService{
		Client: rdb,
	}, nil

}
