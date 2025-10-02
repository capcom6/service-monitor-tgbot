package infrastructure

import (
	"context"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TelegramBotParams struct {
	fx.In

	Config   config.Telegram
	Logger   *zap.Logger
	Monitor  *monitor.MonitorModule
	Messages interface {
		RenderStatus(services []monitor.ServiceStatus) string
	}
}

type TelegramBot struct {
	Config   config.Telegram
	Logger   *zap.Logger
	Monitor  *monitor.MonitorModule
	Messages interface {
		RenderStatus(services []monitor.ServiceStatus) string
	}

	api *tgbotapi.BotAPI
}

func NewTelegramBot(p TelegramBotParams) (*TelegramBot, error) {
	api, err := tgbotapi.NewBotAPI(p.Config.Token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		Config:   p.Config,
		Logger:   p.Logger,
		Monitor:  p.Monitor,
		Messages: p.Messages,
		api:      api,
	}, nil
}

func (b *TelegramBot) SendMessage(text string) error {
	msg := tgbotapi.NewMessage(b.Config.ChatID, text)
	msg.ParseMode = tgbotapi.ModeMarkdownV2
	_, err := b.api.Send(msg)
	return err
}

func (b *TelegramBot) EscapeText(text string) string {
	return tgbotapi.EscapeText(tgbotapi.ModeMarkdownV2, text)
}

func (b *TelegramBot) Listen(ctx context.Context) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "callback_query"}

	go func() {
		<-ctx.Done()
		b.api.StopReceivingUpdates()
	}()

	updates := b.api.GetUpdatesChan(u)

	// Process updates in a separate goroutine
	go func() {
		for update := range updates {
			if update.Message != nil && update.Message.IsCommand() {
				b.handleCommand(ctx, update.Message)
			}
		}
	}()

	return updates, nil
}

func (b *TelegramBot) handleCommand(ctx context.Context, msg *tgbotapi.Message) {
	if msg.Command() == "status" {
		b.HandleStatusCommand(msg.Chat.ID)
	}
}

func (b *TelegramBot) HandleStatusCommand(chatID int64) {
	// Get current service states
	serviceStatuses := b.getCurrentServiceStatuses()

	// Format response using RenderStatus method
	response := b.Messages.RenderStatus(serviceStatuses)

	// Send response
	if err := b.SendMessage(response); err != nil {
		b.Logger.Error("Failed to send status message", zap.Error(err))
	}
}

func (b *TelegramBot) getCurrentServiceStatuses() []monitor.ServiceStatus {
	return b.Monitor.GetCurrentStatuses()
}
