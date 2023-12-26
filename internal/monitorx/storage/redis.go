package storage

import (
	"context"
	"encoding/json"
	"net/url"

	"github.com/capcom6/service-monitor-tgbot/pkg/collections"
	"github.com/redis/go-redis/v9"
)

type redisStorage struct {
	conn *redis.Client
}

func newRedisStorage(u *url.URL) (*redisStorage, error) {
	opts, err := redis.ParseURL(u.String())
	if err != nil {
		return nil, err
	}

	return &redisStorage{
		conn: redis.NewClient(opts),
	}, nil
}

func (s *redisStorage) Select(ctx context.Context) ([]Service, error) {
	vals, err := s.conn.HVals(ctx, "services").Result()
	if err != nil {
		return nil, err
	}

	services, err := collections.Map[string, Service](vals, func(s string) (Service, error) {
		service := Service{}
		return service, json.Unmarshal([]byte(s), &service)
	})

	return services, nil
}
