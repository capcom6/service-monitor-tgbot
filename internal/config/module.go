package config

import (
	"github.com/capcom6/service-monitor-tgbot/internal/bot"
	"github.com/capcom6/service-monitor-tgbot/internal/heartbeat"
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"github.com/capcom6/service-monitor-tgbot/internal/storage"
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
		// httpfx owns the transport: SOCKS5 routing from ProxyURL and the
		// Timeout that bounds dialing the proxy, TLS and reading the response
		// in one deadline. Every other field stays at its inherit-the-defaults
		// zero (pools, idle timeout, root CAs). telegofx receives the
		// configured *http.Client from the graph (httpfx.Module).
		fx.Provide(func(cfg Config) httpfx.Config {
			return httpfx.Config{
				ProxyURL:            cfg.Telegram.ProxyURL,
				Timeout:             cfg.Telegram.Timeout,
				Bypass:              "",
				MaxIdleConns:        0,
				MaxIdleConnsPerHost: 0,
				IdleConnTimeout:     0,
				TLS: httpfx.TLSConfig{
					RootCAFile:          "",
					RootCAPEM:           "",
					RootCAReplaceSystem: false,
				},
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
