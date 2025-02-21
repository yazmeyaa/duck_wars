package redis

import (
	"context"
	"net"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/yazmeyaa/ton_game/internal/config"
)

func NewRedisClient(cfg *config.Config) (*redis.Client, error) {
	addr := net.JoinHostPort(cfg.Redis.Host, cfg.Redis.Port)

	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		DB:       cfg.Redis.DB,
		Password: cfg.Redis.Password,
		Username: cfg.Redis.Username,
	})

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*2)
	defer cancel()
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return rdb, nil
}
