package monitorx

import (
	"context"
	"sync"
	"time"

	"github.com/capcom6/go-infra-fx/logger"
	"github.com/capcom6/service-monitor-tgbot/internal/monitorx/config"
	"github.com/capcom6/service-monitor-tgbot/internal/monitorx/eventbus"
	"github.com/capcom6/service-monitor-tgbot/internal/monitorx/monitor"
	"github.com/capcom6/service-monitor-tgbot/internal/monitorx/storage"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	module := fx.Module(
		"bot",
		logger.Module,
		config.Module,
		monitor.Module,
		storage.Module,
		eventbus.Module,
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger, monitorMod *monitor.MonitorModule, eventbus eventbus.EventBus) error {
			ctx, cancel := context.WithCancel(context.Background())
			wg := &sync.WaitGroup{}
			lc.Append(fx.Hook{
				OnStart: func(ctx context.Context) error {
					logger.Info("Started")
					return nil
				},
				OnStop: func(ctx context.Context) error {
					cancel()
					wg.Wait()
					logger.Info("Stopped")
					return nil
				},
			})

			ch, err := monitorMod.Monitor(ctx)
			if err != nil {
				return err
			}

			wg.Add(1)
			go func() {
				defer wg.Done()
				for v := range ch {
					logger.Debug("probe", zap.Any("state", v))

					ctx, cancel := context.WithTimeout(ctx, 1*time.Second)
					if err := eventbus.Send(ctx, v); err != nil {
						logger.Error("failed to send event", zap.Error(err))
					}
					cancel()
				}
			}()

			return nil
		}),
	)

	fx.New(
		module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
	).Run()
}
