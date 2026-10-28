package config_test

// These tests measure the production Telegram transport path: the TELEGRAM__*
// environment keys through config.Module's provider mapping, and the module
// graph the app boots (httpfx + telegram + telegofx). Everything is
// construction-only - fx.New never starts, so no polling begins and no packet
// leaves the machine: the rejection rows stop inside telegram.Validate, httpfx
// or telego's token gate.
//
// The config is driven through the real env path (t.Setenv, so these tests are
// deliberately not parallel), because the mapping itself is what is measured.

import (
	"errors"
	"testing"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/telegram"
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
// transport: telegofx gets the token and the polling mode, telegram.Config gets
// the proxy and the timeout for the consumer to validate and apply, and the
// app-wide httpfx.Config stays at its inherit-the-defaults zero.
func TestModule_TelegramTransportMapping(t *testing.T) {
	setTelegramEnv(t, mappingToken, acceptedProxy, mappingTimeout)

	var tgCfg telegofx.Config
	var tgTransportCfg telegram.Config
	var httpCfg httpfx.Config

	err := fx.New(
		config.Module(),
		fx.Populate(&tgCfg, &tgTransportCfg, &httpCfg),
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
	if tgTransportCfg.ProxyURL != acceptedProxy {
		t.Errorf("telegram proxy URL = %q, want %q", tgTransportCfg.ProxyURL, acceptedProxy)
	}
	if tgTransportCfg.Timeout != 6*time.Minute {
		t.Errorf("telegram timeout = %s, want %s", tgTransportCfg.Timeout, 6*time.Minute)
	}
	if httpCfg != (httpfx.Config{}) {
		t.Errorf("httpfx config = %+v, want the inherit-the-defaults zero", httpCfg)
	}
}

// TestModule_ShortTimeoutRejected pins startup rejection of a telegram timeout
// that cannot outlive telegofx's fixed 8s getUpdates long-poll hold: the
// telegram module validates before it builds the bot's client, so the graph
// never resolves and no polling cycle can be cancelled mid-hold.
func TestModule_ShortTimeoutRejected(t *testing.T) {
	setTelegramEnv(t, mappingToken, acceptedProxy, "8s")

	var tgOpts []telegofx.Option

	err := fx.New(
		config.Module(),
		httpfx.Module(),
		telegram.Module(),
		fx.Populate(&tgOpts),
		fx.WithLogger(func() fxevent.Logger { return fxevent.NopLogger }),
	).Err()
	if err == nil {
		t.Fatal("graph(timeout=8s) resolved, want a startup failure")
	}
	if !errors.Is(err, telegram.ErrConfigInvalid) {
		t.Errorf("error = %v, want errors.Is telegram.ErrConfigInvalid", err)
	}
}

// TestModule_TelegramGraphGates pins the startup failure order of the wiring the
// app actually boots (config + httpfx + telegram + telegofx with router): the
// bot's client is built as a dependency of telegofx.New via []telegofx.Option,
// so a rejected proxy surfaces first, and only an accepted proxy lets telego's
// token gate run.
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
				telegram.Module(),
				telegofx.Module(telegofx.WithRouter()),
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
