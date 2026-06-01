package config

import "time"

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	// Services []Service `yaml:"services"`
}

type Telegram struct {
	Token      string           `yaml:"token" envconfig:"TELEGRAM__TOKEN" validate:"required"`
	ChatID     int64            `yaml:"chatId" envconfig:"TELEGRAM__CHAT_ID"`
	WebhookURL string           `yaml:"webhookUrl" envconfig:"TELEGRAM__WEBHOOK_URL" validate:"required"`
	Debug      bool             `yaml:"debug" envconfig:"TELEGRAM__DEBUG"`
	ProxyURL   string           `yaml:"proxyUrl" envconfig:"TELEGRAM__PROXY_URL"`
	Timeout    time.Duration    `yaml:"timeout" envconfig:"TELEGRAM__TIMEOUT"`
	Messages   TelegramMessages `yaml:"messages"`
}

type TelegramMessages map[string]string
