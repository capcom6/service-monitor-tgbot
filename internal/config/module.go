package config

import (
	"github.com/capcom6/go-infra-fx/config"
	"github.com/capcom6/go-infra-fx/fxutil"
	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		fxutil.WithNamedLogger("config"),
		fx.Provide(func() (Config, error) {
			return defaultConfig, config.LoadConfig(&defaultConfig)
		}),
		fx.Provide(func(cfg Config) telegram.Config {
			return telegram.Config{
				Token: cfg.Telegram.Token,
			}
		}),
		fx.Provide(func(cfg Config) messages.Config {
			return messages.Config{
				Templates: cfg.Telegram.Messages,
			}
		}),
		fx.Provide(func(cfg Config) bot.Config {
			return bot.Config{
				ChatID: cfg.Telegram.ChatID,
			}
		}),
	)
}
