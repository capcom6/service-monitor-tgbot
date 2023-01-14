package infrastructure

import (
	"fmt"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
	"github.com/go-redis/redis/v9"
)

func redisConnect(cfg config.Redis) (*redis.Client, error) {
	return redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		Password: cfg.Password, // no password set
		DB:       cfg.DB,       // use default DB
	}), nil
}
