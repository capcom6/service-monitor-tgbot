package telegram

import (
	"context"
)

type Command struct {
	From int64
	Name string
	Args string
}

type CommandHandler func(ctx context.Context, cmd Command)
