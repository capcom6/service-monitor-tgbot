package messages

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"messages",
		logger.WithNamedLogger("messages"),
		fx.Provide(NewService),
	)
}
