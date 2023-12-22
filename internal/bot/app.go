package bot

import (
	"context"
	"sync"

	"github.com/capcom6/go-infra-fx/logger"
	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/infrastructure"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Module = fx.Module(
	"bot",
	logger.Module,
	config.Module,
	infrastructure.Module,
	monitor.Module,
	fx.Provide(NewMessages),
)

func Run() {
	fx.New(
		Module,
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger, bot *infrastructure.TelegramBot, monitorMod *monitor.MonitorModule, messages *Messages) error {
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
				for v := range ch {
					logger.Debug("probe", zap.String("name", v.Name), zap.String("state", string(v.State)), zap.Error(v.Error))

					msg := ""
					if v.State == monitor.ServiceOffline {
						context := OfflineContext{
							OnlineContext: OnlineContext{
								Name: bot.EscapeText(v.Name),
							},
							Error: bot.EscapeText(v.Error.Error()),
						}
						msg, err = messages.Render(TemplateOffline, context)
					} else {
						context := OnlineContext{
							Name: bot.EscapeText(v.Name),
						}
						msg, err = messages.Render(TemplateOnline, context)
					}

					if err != nil {
						logger.Error("can't render template", zap.Error(err))
						continue
					}

					if err := bot.SendMessage(msg); err != nil {
						logger.Error("can't send message", zap.Error(err))
					}
				}
			}()

			return nil
		}),
	).Run()
}
