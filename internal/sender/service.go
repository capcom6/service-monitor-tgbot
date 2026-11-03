// Package sender delivers pre-rendered messages to Telegram chats.
package sender

import (
	"context"
	"fmt"

	"github.com/go-core-fx/telegofx"
	"github.com/mymmrac/telego"
)

// parseMode is the Markdown dialect messages are sent in. It is the mode
// EscapeText escapes for, so the two must stay in sync.
const parseMode = telego.ModeMarkdownV2

// Service sends messages through the bot the telegofx module runs.
type Service struct {
	bot *telego.Bot
}

// New creates a message sender on top of the telegofx bot.
func New(bot *telegofx.Bot) *Service {
	return &Service{bot: bot.Bot}
}

// Send sends text to chatID in the package's parse mode and returns the Bot
// API failure, if any.
func (s *Service) Send(ctx context.Context, chatID int64, text string) error {
	_, err := s.bot.SendMessage(ctx, &telego.SendMessageParams{
		ChatID:    telego.ChatID{ID: chatID, Username: ""},
		Text:      text,
		ParseMode: parseMode,
	})
	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}
