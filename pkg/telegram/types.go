package telegram

import (
	"context"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type CommandHandler interface {
	HandleCommand(ctx context.Context, msg *tg.Message)
}
