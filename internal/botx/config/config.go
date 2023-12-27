package config

type Config struct {
	Telegram Telegram `yaml:"telegram"`
	EventBus EventBus `yaml:"eventBus"`
}

type Telegram struct {
	Token      string           `yaml:"token" envconfig:"TELEGRAM__TOKEN" validate:"required"`
	ChatID     int64            `yaml:"chatId" envconfig:"TELEGRAM__CHAT_ID"`
	WebhookURL string           `yaml:"webhookUrl" envconfig:"TELEGRAM__WEBHOOK_URL" validate:"required"`
	Debug      bool             `yaml:"debug" envconfig:"TELEGRAM__DEBUG"`
	Messages   TelegramMessages `yaml:"messages"`
}

type TelegramMessages map[string]string

type EventBus struct {
	DSN string `yaml:"dsn" envconfig:"EVENTBUS__DSN" validate:"required"`
}
