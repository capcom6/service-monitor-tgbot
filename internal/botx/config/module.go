package config

import (
	"github.com/capcom6/go-infra-fx/config"
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
)
