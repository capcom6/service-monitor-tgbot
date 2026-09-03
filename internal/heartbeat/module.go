package heartbeat

import (
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

// Module creates and returns an FX module for the heartbeat package.
func Module() fx.Option {
	return fx.Module(
		"heartbeat",
		logger.WithNamedLogger("heartbeat"),
		fx.Provide(New),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
	)
}
