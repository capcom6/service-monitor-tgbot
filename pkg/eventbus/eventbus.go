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
	Send(ctx context.Context, event string) error
	Receive(ctx context.Context) (<-chan string, error)
}

type Params struct {
	fx.In

	Config Config
	Logger *zap.Logger
}

var Module = fx.Module(
	"eventbus",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("eventbus")
	}),
	fx.Provide(func(p Params) (EventBus, error) {
		return New(p.Config, p.Logger)
	}),
)

func New(cfg Config, logger *zap.Logger) (EventBus, error) {
	u, err := url.Parse(cfg.DSN)
	if err != nil {
		return nil, fmt.Errorf("invalid dsn: %w", err)
	}

	if u.Scheme == "redis" || u.Scheme == "rediss" {
		return newRedisBus(u, logger)
	}

	return nil, errors.New("unknown scheme")
}
