package telegram_test

import (
	"testing"
	"time"

	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
)

func TestConfig_LongPollTimeout(t *testing.T) {
	tests := []struct {
		name    string
		timeout time.Duration
		want    int
	}{
		{
			name:    "zero timeout results in short polling",
			timeout: 0,
			want:    0,
		},
		{
			name:    "timeout below margin results in short polling",
			timeout: 3 * time.Second,
			want:    0,
		},
		{
			name:    "default timeout leaves headroom below client timeout",
			timeout: time.Minute,
			want:    55,
		},
		{
			name:    "deployed timeout of 30s stays below client timeout",
			timeout: 30 * time.Second,
			want:    25,
		},
		{
			name:    "timeout above max caps long poll at telegram limit",
			timeout: 90 * time.Second,
			want:    60,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := telegram.Config{Timeout: tt.timeout}
			if got := cfg.LongPollTimeout(); got != tt.want {
				t.Fatalf("LongPollTimeout() = %d, want %d", got, tt.want)
			}
		})
	}
}
