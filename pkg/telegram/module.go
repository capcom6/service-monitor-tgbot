package telegram

import (
	"context"
	"sync"

	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"telegram",
		fxutil.WithNamedLogger("telegram"),
		fx.Provide(NewBot),
		fx.Invoke(func(bot *Bot, logger *zap.Logger, lifecycle fx.Lifecycle, shutdowner fx.Shutdowner) {
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}

			lifecycle.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("Starting Telegram bot")
					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := bot.Listen(ctx); err != nil {
							logger.Error("Failed to run bot", zap.Error(err))
							if shutdownErr := shutdowner.Shutdown(fx.ExitCode(1)); shutdownErr != nil {
								logger.Error("Failed to trigger shutdown", zap.Error(shutdownErr))
							}
						}
					}()
					logger.Info("Telegram bot started")

					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("Stopping Telegram bot")
					cancel()
					bot.Close()
					wg.Wait()
					logger.Info("Telegram bot stopped")

					return nil
				},
			})
		}),
	)
}
