package telegram

import (
	"github.com/capcom6/service-monitor-tgbot/internal/botx/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type TelegramBotParams struct {
	fx.In

	Config config.Telegram
	Logger *zap.Logger
}

type TelegramBot struct {
	Config config.Telegram

	api *tgbotapi.BotAPI
}

func NewTelegramBot(p TelegramBotParams) (*TelegramBot, error) {
	api, err := tgbotapi.NewBotAPI(p.Config.Token)
	if err != nil {
		return nil, err
	}

	return &TelegramBot{
		Config: p.Config,
		api:    api,
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
