package internal

import (
	"github.com/capcom6/go-infra-fx/logger"
	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func Run() {
	fx.New(
		fx.Module(
			"app",
			logger.Module,
			config.Module(),
			messages.Module(),
			telegram.Module(),
			monitor.Module(),
			storage.Module(),
			bot.Module(),
		),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			logOption := fxevent.ZapLogger{Logger: logger}
			logOption.UseLogLevel(zapcore.DebugLevel)
			return &logOption
		}),
	).Run()
}
