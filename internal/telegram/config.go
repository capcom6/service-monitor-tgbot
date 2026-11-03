package telegram

import (
	"fmt"
	"time"
)

type Config struct {
	// ProxyURL is the SOCKS5 proxy the bot's API traffic dials through.
	ProxyURL string
	// Timeout bounds every outgoing Bot API request end to end.
	Timeout time.Duration
}

const minTelegramTimeout = 8 * time.Second

// Validate checks the config where it is consumed, before it builds the
// bot's HTTP client: the timeout must outlive telego's 8s getUpdates
// long-poll hold (ErrConfigInvalid wraps the detail).
func (c Config) Validate() error {
	if c.Timeout <= minTelegramTimeout {
		return fmt.Errorf(
			"%w: telegram timeout %v must be greater than the 8s getUpdates long-poll hold",
			ErrConfigInvalid, c.Timeout,
		)
	}

	return nil
}
