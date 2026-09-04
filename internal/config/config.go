package config

import (
	"fmt"
	"os"
	"time"

	"github.com/go-core-fx/config"
)

const defaultHeartbeatInterval = 6 * time.Hour

type http struct {
	Address     string   `koanf:"address"`
	ProxyHeader string   `koanf:"proxy_header"`
	Proxies     []string `koanf:"proxies"`

	OpenAPI openAPIConfig `koanf:"openapi"`
}

type openAPIConfig struct {
	Enabled    bool   `koanf:"enabled"`
	PublicHost string `koanf:"public_host"`
	PublicPath string `koanf:"public_path"`
}

type telegramConfig struct {
	Token    string           `koanf:"token"`
	ProxyURL string           `koanf:"proxy_url"`
	Timeout  time.Duration    `koanf:"timeout"`
	ChatID   int64            `koanf:"chatId"`
	Debug    bool             `koanf:"debug"`
	Messages TelegramMessages `koanf:"messages"`
}

type TelegramMessages map[string]string

type storageConfig struct {
	DSN string `koanf:"dsn"`
}

type heartbeatConfig struct {
	Enabled  bool          `koanf:"enabled"`
	Interval time.Duration `koanf:"interval"`
	ChatID   int64         `koanf:"chatId"`
}

type Config struct {
	HTTP      http            `koanf:"http"`
	Telegram  telegramConfig  `koanf:"telegram"`
	Storage   storageConfig   `koanf:"storage"`
	Heartbeat heartbeatConfig `koanf:"heartbeat"`
}

func Default() Config {
	return Config{
		HTTP: http{
			Address:     "127.0.0.1:3000",
			ProxyHeader: "X-Forwarded-For",
			Proxies:     []string{},
			OpenAPI: openAPIConfig{
				Enabled:    true,
				PublicHost: "",
				PublicPath: "",
			},
		},
		Telegram: telegramConfig{
			Token:    "",
			ProxyURL: "",
			Timeout:  time.Minute,
			ChatID:   0,
			Debug:    false,
			Messages: TelegramMessages{},
		},
		Storage: storageConfig{
			DSN: "file://" + os.Getenv("CONFIG_PATH"),
		},
		Heartbeat: heartbeatConfig{
			Enabled:  false,
			Interval: defaultHeartbeatInterval,
			ChatID:   0,
		},
	}
}

func New() (Config, error) {
	cfg := Default()

	options := []config.Option{}
	if yamlPath := os.Getenv("CONFIG_PATH"); yamlPath != "" {
		options = append(options, config.WithLocalYAML(yamlPath))
	}

	if err := config.Load(&cfg, options...); err != nil {
		return Config{}, fmt.Errorf("failed to load config: %w", err)
	}

	return cfg, nil
}
