package internal

import (
	"context"

	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/internal/server"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/healthfx"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Run(version healthfx.Version) {
	fx.New(
		// CORE MODULES
		logger.Module(),
		logger.WithFxDefaultLogger(),
		// badgerfx.Module(),
		// bunfx.Module(),
		// cachefx.Module(),
		fiberfx.Module(),
		// gocqlfx.Module(),
		// gocqlxfx.Module(),
		// sqlfx.Module(),
		// goosefx.Module(),
		// gormfx.Module(),
		healthfx.Module(),
		// openrouterfx.Module(),
		// redisfx.Module(),
		// sqlxfx.Module(),
		// telegofx.Module(true),
		// validatorfx.Module(),
		// watermillfx.Module(),
		//
		// APP MODULES
		config.Module(),
		// db.Module(),
		server.Module(),
		// bot.Module(),
		telegram.Module(),
		storage.Module(),
		bot.Module(),
		//
		// BUSINESS MODULES
		fx.Supply(version),
		messages.Module(),
		monitor.Module(),
		//
		fx.Invoke(func(lc fx.Lifecycle, logger *zap.Logger) {
			lc.Append(fx.Hook{
				OnStart: func(_ context.Context) error {
					logger.Info("app started")
					return nil
				},
				OnStop: func(_ context.Context) error {
					logger.Info("app stopped")
					return nil
				},
			})
		}),
	).Run()

	// fx.New(
	// 	fx.Module(
	// 		"app",
	// 		logger.Module,
	// 		config.Module(),
	// 		messages.Module(),
	// 		telegram.Module(),
	// 		monitor.Module(),
	// 		storage.Module(),
	// 		bot.Module(),
	// 	),
	// 	fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
	// 		logOption := fxevent.ZapLogger{Logger: logger}
	// 		logOption.UseLogLevel(zapcore.DebugLevel)
	// 		return &logOption
	// 	}),
	// ).Run()
}
