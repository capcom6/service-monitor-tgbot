package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/redis/go-redis/v9"
)

const defaultRedisKey = "service-monitor:services"

type RedisStorage struct {
	key string

	client *redis.Client
}

func newRedisStorage(u *url.URL) (Storage, error) {
	opts, err := redisOptionsFromURL(u)
	if err != nil {
		return nil, err
	}

	client := redis.NewClient(opts)

	key := u.Query().Get("key")
	if key == "" {
		key = defaultRedisKey
	}

	return &RedisStorage{
		client: client,
		key:    key,
	}, nil
}

func redisOptionsFromURL(u *url.URL) (*redis.Options, error) {
	q := u.Query()
	q.Del("key")

	uCopy := *u
	uCopy.RawQuery = q.Encode()

	opts, err := redis.ParseURL(uCopy.String())
	if err != nil {
		return nil, fmt.Errorf("%w: %w", ErrRedisInvalidDSN, err)
	}

	return opts, nil
}

func (s *RedisStorage) Load() ([]MonitoredService, error) {
	data, err := s.client.Get(context.Background(), s.key).Bytes()
	if err != nil {
		return nil, fmt.Errorf("failed to get services from redis: %w", err)
	}

	var services []MonitoredService
	if jsonErr := json.Unmarshal(data, &services); jsonErr != nil {
		return nil, fmt.Errorf("failed to unmarshal services: %w", jsonErr)
	}

	return services, nil
}

func (s *RedisStorage) Close() error {
	if s.client != nil {
		if closeErr := s.client.Close(); closeErr != nil {
			return fmt.Errorf("failed to close redis client: %w", closeErr)
		}
	}
	return nil
}
