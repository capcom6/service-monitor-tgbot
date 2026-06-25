package storage

import (
	"context"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"storage",
		logger.WithNamedLogger("storage"),
		fx.Provide(func(cfg Config, lc fx.Lifecycle) (Storage, error) {
			s, err := NewFromDSN(cfg.DSN)
			if err != nil {
				return nil, err
			}

			lc.Append(fx.Hook{
				OnStart: nil,
				OnStop: func(_ context.Context) error {
					return s.Close()
				},
			})

			return NewService(s), nil
		}),
	)
}
