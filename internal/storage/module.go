package storage

import (
	"errors"
	"fmt"
	"net/url"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type StorageParams struct {
	fx.In

	Config Config
	Logger *zap.Logger
}

var Module = fx.Module(
	"storage",
	fx.Provide(func(p StorageParams) (Storage, error) {
		s, err := newStorage(p)
		if err != nil {
			return nil, err
		}

		return NewStorageService(s), nil
	}),
)

func newStorage(p StorageParams) (Storage, error) {
	u, err := url.Parse(p.Config.DSN)
	if err != nil {
		return nil, fmt.Errorf("invalid dsn: %w", err)
	}

	if u.Scheme == "file" {
		return newYamlStorage(u)
	}
	if u.Scheme == "redis" || u.Scheme == "rediss" {
		return newRedisStorage(u)
	}
	if u.Scheme == "memory" {
		return newMemoryStorage(u)
	}

	return nil, errors.New("unknown scheme")
}
