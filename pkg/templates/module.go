package templates

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"templates",
		logger.WithNamedLogger("templates"),
		fx.Provide(NewService),
	)
}
