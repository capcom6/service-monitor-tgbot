package telegram

import (
	"time"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

const (
	maxLongPollTimeoutSeconds = 60
	longPollTimeoutMargin     = 5 * time.Second
)

type Config struct {
	Token     string
	ParseMode string
	ProxyURL  string
	Timeout   time.Duration
}

// LongPollTimeout returns the long polling timeout in seconds for getUpdates.
// It is always kept below the HTTP client timeout so an idle long poll never
// aborts with a context deadline error.
func (c Config) LongPollTimeout() int {
	clientTimeout := c.Timeout - longPollTimeoutMargin
	poll := int(clientTimeout.Seconds())
	if poll <= 0 {
		return 0
	}
	if poll > maxLongPollTimeoutSeconds {
		return maxLongPollTimeoutSeconds
	}
	return poll
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
