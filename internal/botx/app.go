package botx

import (
	"context"
	"sync"

	"github.com/capcom6/go-infra-fx/logger"
	"github.com/capcom6/service-monitor-tgbot/internal/botx/config"
	"github.com/capcom6/service-monitor-tgbot/internal/botx/telegram"
	"github.com/capcom6/service-monitor-tgbot/pkg/eventbus"
	"github.com/capcom6/service-monitor-tgbot/pkg/events"
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
		telegram.Module,
		eventbus.Module,
		fx.Provide(NewMessages),
		fx.Invoke(start),
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

func start(eventbus eventbus.EventBus, telegram *telegram.TelegramBot, messages *Messages, logger *zap.Logger, lc fx.Lifecycle) error {
	ctx, cancel := context.WithCancel(context.Background())

	wg := &sync.WaitGroup{}

	lc.Append(fx.Hook{
		OnStop: func(ctx context.Context) error {
			cancel()
			wg.Wait()
			return nil
		},
	})

	ch, err := eventbus.Receive(ctx)
	if err != nil {
		return err
	}

	wg.Add(1)
	go func() {
		defer wg.Done()

		logger.Debug("start receive events")
		defer logger.Debug("stop receive events")

		event := events.Event[events.ServiceStateChanged]{}
		for payload := range ch {
			if err := event.Decode(payload); err != nil {
				continue
			}
			if event.Name != events.EventNameServiceStateChanged {
				continue
			}

			v := event.Payload

			msg := ""
			if v.State == events.ServiceStateOffline {
				context := OfflineContext{
					OnlineContext: OnlineContext{
						Name: telegram.EscapeText(v.Name),
					},
					Error: telegram.EscapeText(v.Error),
				}
				msg, err = messages.Render(TemplateOffline, context)
			} else {
				context := OnlineContext{
					Name: telegram.EscapeText(v.Name),
				}
				msg, err = messages.Render(TemplateOnline, context)
			}

			if err != nil {
				logger.Error("can't render template", zap.Error(err))
				continue
			}

			if err := telegram.SendMessage(msg); err != nil {
				logger.Error("can't send message", zap.Error(err))
			}
		}
	}()

	return nil
}
