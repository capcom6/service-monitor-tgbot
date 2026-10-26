package config_test

// These tests measure the production Telegram transport path: the TELEGRAM__*
// environment keys through config.Module's provider mapping, and the module
// graph the app boots (httpfx + telegofx + the telegram sender module).
// Everything is construction-only - fx.New never starts, so no polling begins
// and no packet leaves the machine: the rejection rows stop inside httpfx or
// telego's token gate.
//
// The config is driven through the real env path (t.Setenv, so these tests are
// deliberately not parallel), because the mapping itself is what is measured.

import (
	"errors"
	"testing"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"github.com/go-core-fx/httpfx"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
)

const (
	// mappingToken is an arbitrary token value: the mapping test never builds a bot.
	mappingToken = "token-from-env"
	// acceptedProxy is a SOCKS5 URL httpfx accepts at construction.
	acceptedProxy = "socks5://127.0.0.1:1080"
	// rejectedProxy is an HTTP proxy URL, which httpfx rejects: only socks5:// and
	// socks5h:// schemes are registered.
	rejectedProxy = "http://127.0.0.1:8080"
	// mappingTimeout is the raw TELEGRAM__TIMEOUT value; the mapped duration is 6m.
	mappingTimeout = "6m"
)

// setTelegramEnv pins a complete TELEGRAM__* set (plus an empty CONFIG_PATH so a
// developer's local YAML file never leaks into the measurement).
func setTelegramEnv(t *testing.T, token, proxyURL, timeout string) {
	t.Helper()

	t.Setenv("CONFIG_PATH", "")
	t.Setenv("TELEGRAM__TOKEN", token)
	t.Setenv("TELEGRAM__PROXY_URL", proxyURL)
	t.Setenv("TELEGRAM__TIMEOUT", timeout)
}

// TestModule_TelegramTransportMapping pins what config.Module hands the
// transport: telegofx gets the token and the polling mode, httpfx gets the
// proxy and the timeout, and the remaining httpfx fields stay at their
// inherit-the-defaults zero.
func TestModule_TelegramTransportMapping(t *testing.T) {
	setTelegramEnv(t, mappingToken, acceptedProxy, mappingTimeout)

	var tgCfg telegofx.Config
	var httpCfg httpfx.Config

	err := fx.New(
		config.Module(),
		fx.Populate(&tgCfg, &httpCfg),
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
	).Err()
	if err != nil {
		t.Fatalf("config graph must resolve: %v", err)
	}

	if tgCfg.Token != mappingToken {
		t.Errorf("telegofx token = %q, want %q", tgCfg.Token, mappingToken)
	}
	if tgCfg.Mode != telegofx.ModePolling {
		t.Errorf("telegofx mode = %q, want %q", tgCfg.Mode, telegofx.ModePolling)
	}
	if httpCfg.ProxyURL != acceptedProxy {
		t.Errorf("httpfx proxy URL = %q, want %q", httpCfg.ProxyURL, acceptedProxy)
	}
	if httpCfg.Timeout != 6*time.Minute {
		t.Errorf("httpfx timeout = %s, want %s", httpCfg.Timeout, 6*time.Minute)
	}
	if httpCfg.TLS != (httpfx.TLSConfig{}) {
		t.Errorf("httpfx TLS = %+v, want the inherit-the-defaults zero", httpCfg.TLS)
	}
}

// TestModule_TelegramGraphGates pins the startup failure order of the wiring the
// app actually boots (config + httpfx + telegofx with router + sender): the
// httpfx client is built as a dependency of telegofx.New, so a rejected proxy
// surfaces first, and only an accepted proxy lets telego's token gate run. The
// empty-token row also proves sender.Module supplies the []telego.BotOption
// telegofx.New requires: without it dig would fail with a missing type instead.
func TestModule_TelegramGraphGates(t *testing.T) {
	tests := []struct {
		name     string
		proxyURL string
		wantErr  error
	}{
		{
			name:     "rejected proxy fails in httpfx before the token gate",
			proxyURL: rejectedProxy,
			wantErr:  httpfx.ErrInvalidProxyURL,
		},
		{name: "accepted proxy reaches telego's token gate", proxyURL: acceptedProxy, wantErr: telego.ErrEmptyToken},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			setTelegramEnv(t, "", tt.proxyURL, "1m")

			err := fx.New(
				logger.Module(),
				httpfx.Module(),
				config.Module(),
				telegofx.Module(telegofx.WithRouter()),
				sender.Module(),
				fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
			).Err()
			if err == nil {
				t.Fatalf("graph(proxy=%q) resolved, want a startup failure", tt.proxyURL)
			}
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("graph(proxy=%q) error %v, want errors.Is %v", tt.proxyURL, err, tt.wantErr)
			}
		})
	}
}
