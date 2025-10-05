package telegram

import (
	"context"

	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"telegram",
		fxutil.WithNamedLogger("telegram"),
		fx.Provide(NewBot),
		fx.Invoke(func(bot *Bot, logger *zap.Logger, lifecycle fx.Lifecycle) {
			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("Starting Telegram bot")
					go bot.Listen()
					logger.Info("Telegram bot started")

					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("Stopping Telegram bot")
					bot.Close()
					logger.Info("Telegram bot stopped")

					return nil
				},
			})
		}),
	)
}
