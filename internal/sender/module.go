package sender

import (
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

// Module provides the message Service.
func Module() fx.Option {
	return fx.Module(
		"sender",
		logger.WithNamedLogger("sender"),
		// The transport is telegram.Module's private *http.Client handed to
		// telegofx through []telegofx.Option (SOCKS5 routing and timeout), so
		// no bot option is configured here.
		fx.Provide(New),
	)
}
