package config

import (
	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/logger"
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
		fx.Provide(func(cfg Config) telegram.Config {
			return telegram.Config{
				Token:     cfg.Telegram.Token,
				ParseMode: "MarkdownV2",
				ProxyURL:  cfg.Telegram.ProxyURL,
				Timeout:   cfg.Telegram.Timeout,
			}
		}),
		fx.Provide(func(cfg Config) messages.Config {
			return messages.Config{
				Templates: cfg.Telegram.Messages,
				EscapeFn:  nil,
			}
		}),
		fx.Provide(func(cfg Config) bot.Config {
			return bot.Config{
				ChatID: cfg.Telegram.ChatID,
			}
		}),
		fx.Provide(func(cfg Config) storage.Config {
			return storage.Config{
				DSN: cfg.Storage.DSN,
			}
		}),
	)
}
