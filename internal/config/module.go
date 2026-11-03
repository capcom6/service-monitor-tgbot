package config

import (
	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/heartbeat"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/capcom6/service-monitor-tgbot/internal/telegram"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/httpfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/samber/lo"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"config",
		logger.WithNamedLogger("config"),
		fx.Provide(New, fx.Private),
		fx.Provide(
			func(cfg Config) fiberfx.Config {
				return fiberfx.Config{
					Address:     cfg.HTTP.Address,
					ProxyHeader: cfg.HTTP.ProxyHeader,
					Proxies:     cfg.HTTP.Proxies,
				}
			},
			func(cfg Config) openapi.Config {
				return openapi.Config{
					Enabled:    cfg.HTTP.OpenAPI.Enabled,
					PublicHost: cfg.HTTP.OpenAPI.PublicHost,
					PublicPath: cfg.HTTP.OpenAPI.PublicPath,
				}
			},
		),
		fx.Provide(func(cfg Config) telegofx.Config {
			return telegofx.Config{
				Token: cfg.Telegram.Token,
				Mode:  telegofx.ModePolling,
			}
		}),
		// The app-wide httpfx client stays at its inherit-the-defaults zero:
		// it carries no transport policy of its own. The bot's proxy and
		// timeout map to telegram.Config instead, and telegram.Module
		// validates them and overrides the bot's client through the
		// httpfx.Factory this config seeds.
		fx.Provide(func() httpfx.Config {
			return httpfx.Config{}
		}),
		fx.Provide(func(cfg Config) telegram.Config {
			return telegram.Config{
				ProxyURL: cfg.Telegram.ProxyURL,
				Timeout:  cfg.Telegram.Timeout,
			}
		}),
		fx.Provide(func(cfg Config) messages.Config {
			return messages.Config{
				Templates: cfg.Telegram.Messages,
				EscapeFn:  sender.EscapeText,
			}
		}),
		fx.Provide(func(cfg Config) bot.Config {
			return bot.Config{
				ChatID: cfg.Telegram.ChatID,
			}
		}),
		fx.Provide(func(cfg Config) heartbeat.Config {
			return heartbeat.Config{
				Enabled:  cfg.Heartbeat.Enabled,
				Interval: cfg.Heartbeat.Interval,
				ChatID:   lo.CoalesceOrEmpty(cfg.Heartbeat.ChatID, cfg.Telegram.ChatID),
			}
		}),
		fx.Provide(func(cfg Config) storage.Config {
			return storage.Config{
				DSN: cfg.Storage.DSN,
			}
		}),
	)
}
