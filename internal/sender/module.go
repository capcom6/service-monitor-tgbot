package sender

import (
	"github.com/go-core-fx/logger"
	"github.com/mymmrac/telego"
	"go.uber.org/fx"
)

// Module provides the message Service and the bot options telegofx reads from
// the graph.
func Module() fx.Option {
	return fx.Module(
		"sender",
		logger.WithNamedLogger("sender"),
		// telegofx.New takes []telego.BotOption from the graph. The transport is
		// httpfx's *http.Client (SOCKS5 routing and timeout, provided by
		// httpfx.Module) and telegofx appends the zap logger itself, so no bot
		// option is configured here.
		fx.Provide(func() []telego.BotOption {
			return []telego.BotOption{}
		}),
		fx.Provide(New),
	)
}
