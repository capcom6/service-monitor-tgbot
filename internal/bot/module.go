package bot

import (
	"context"
	"sync"

	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		fxutil.WithNamedLogger("bot"),
		fx.Provide(NewService),
		fx.Invoke(func(s *Service, logger *zap.Logger, lc fx.Lifecycle, sh fx.Shutdowner) {
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}

			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("Starting bot")
					wg.Add(1)
					go func() {
						defer wg.Done()
						if err := s.Run(ctx); err != nil {
							logger.Error("Failed to run bot", zap.Error(err))
							if shutdownErr := sh.Shutdown(fx.ExitCode(1)); shutdownErr != nil {
								logger.Error("Failed to trigger shutdown", zap.Error(shutdownErr))
							}
						}
					}()
					logger.Info("Bot started")

					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("Stopping bot")
					cancel()
					wg.Wait()
					logger.Info("Bot stopped")

					return nil
				},
			})
		}),
	)
}
