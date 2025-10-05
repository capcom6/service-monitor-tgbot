package telegram

import tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"

type Config struct {
	Token     string
	ParseMode string
}

func (c Config) Validate() error {
	if c.Token == "" {
		return ErrTokenIsEmpty
	}

	return nil
}

func (c Config) ApplyDefaults() Config {
	if c.ParseMode == "" {
		c.ParseMode = tg.ModeMarkdownV2
	}

	return c
}
