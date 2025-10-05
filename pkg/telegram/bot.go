package telegram

import (
	"context"
	"fmt"

	tg "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/zap"
)

type Bot struct {
	cfg Config

	api *tg.BotAPI

	handlers map[string]CommandHandler

	logger *zap.Logger
}

func NewBot(config Config, logger *zap.Logger) (*Bot, error) {
	config = config.ApplyDefaults()

	if err := config.Validate(); err != nil {
		return nil, err
	}

	api, err := tg.NewBotAPI(config.Token)
	if err != nil {
		return nil, fmt.Errorf("failed to create bot: %w", err)
	}

	return &Bot{
		cfg: config,

		api: api,

		handlers: make(map[string]CommandHandler),

		logger: logger,
	}, nil
}

func (b *Bot) AddHandler(name string, handler CommandHandler) {
	b.handlers[name] = handler
}

func (b *Bot) SendMessage(chatID int64, text string) (int, error) {
	msg := tg.NewMessage(chatID, text)
	msg.ParseMode = b.cfg.ParseMode

	res, err := b.api.Send(msg)
	if err != nil {
		return 0, fmt.Errorf("failed to send message: %w", err)
	}

	return res.MessageID, nil
}

func (b *Bot) EscapeText(text string) string {
	return tg.EscapeText(b.cfg.ParseMode, text)
}

func (b *Bot) Listen(ctx context.Context) error {
	u := tg.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "callback_query"}

	updates := b.api.GetUpdatesChan(u)
	for update := range updates {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		if update.Message != nil && update.Message.IsCommand() {
			b.handleCommand(ctx, update.Message)
		}
	}

	return nil
}

func (b *Bot) Close() {
	b.api.StopReceivingUpdates()
}

func (b *Bot) handleCommand(ctx context.Context, msg *tg.Message) {
	defer func() {
		if r := recover(); r != nil {
			b.logger.Error("Command handler panicked",
				zap.String("command", msg.Command()),
				zap.Any("panic", r))
		}
	}()

	if handler, ok := b.handlers[msg.Command()]; ok {
		handler(ctx, Command{
			From: msg.From.ID,
			Name: msg.Command(),
			Args: msg.CommandArguments(),
		})
		return
	}

	b.logger.Info("Unknown command", zap.String("command", msg.Command()), zap.Int64("chat_id", msg.Chat.ID))
}
