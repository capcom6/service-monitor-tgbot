package heartbeat

import "time"

// Config holds the configuration for the heartbeat module.
type Config struct {
	Enabled  bool
	Interval time.Duration
	ChatID   int64
}
