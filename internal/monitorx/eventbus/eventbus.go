package eventbus

import (
	"context"
	"errors"
	"fmt"
	"net/url"

	"go.uber.org/fx"
	"go.uber.org/zap"
)

type EventBus interface {
	Send(ctx context.Context, event interface{}) error
}

type Params struct {
	fx.In

	Config Config
	Logger *zap.Logger
}

var Module = fx.Module(
	"storage",
	fx.Provide(func(p Params) (EventBus, error) {
		return New(p.Config)
	}),
)

func New(cfg Config) (EventBus, error) {
	u, err := url.Parse(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("invalid dsn: %w", err)
	}

	if u.Scheme == "redis" || u.Scheme == "rediss" {
		return newRedisBus(u)
	}

	return nil, errors.New("unknown scheme")
}
