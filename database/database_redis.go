package database

import (
	"context"

	goredis "github.com/redis/go-redis/v9"
)

type RedisService struct {
	Client *goredis.Client
}

func NewRedisService() *RedisService {
	client := goredis.NewClient(&goredis.Options{
		Addr: "localhost:6379",
	})

	return &RedisService{
		Client: client,
	}
}

func (r *RedisService) Ping(ctx context.Context) error {
	return r.Client.Ping(ctx).Err()
}
