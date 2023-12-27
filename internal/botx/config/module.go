package config

import (
	"github.com/capcom6/go-infra-fx/config"
	"github.com/capcom6/service-monitor-tgbot/pkg/eventbus"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"appconfig",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("appconfig")
	}),
	fx.Provide(func() (Config, error) {
		return defaultConfig, config.LoadConfig(&defaultConfig)
	}),
	fx.Provide(func(cfg Config) Telegram {
		return cfg.Telegram
	}),
	fx.Provide(func(cfg Config) TelegramMessages {
		return cfg.Telegram.Messages
	}),
	fx.Provide(func(cfg Config) eventbus.Config {
		return eventbus.Config{
			DSN: cfg.EventBus.DSN,
		}
	}),
)
