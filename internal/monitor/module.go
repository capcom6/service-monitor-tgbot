package monitor

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"monitor",
		logger.WithNamedLogger("monitor"),
		fx.Provide(NewService),
	)
}
