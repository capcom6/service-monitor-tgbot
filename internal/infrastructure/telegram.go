package infrastructure

import (
	"context"
	"sync"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type TelegramBot struct {
	Config config.Telegram

	apiOnce sync.Once
	api     *tgbotapi.BotAPI
}

func NewTelegramBot(cfg config.Telegram) *TelegramBot {
	return &TelegramBot{
		Config:  cfg,
		apiOnce: sync.Once{},
	}
}

func (b *TelegramBot) Api() (api *tgbotapi.BotAPI, err error) {
	b.apiOnce.Do(func() {
		b.api, err = tgbotapi.NewBotAPI(b.Config.Token)
		if err == nil {
			b.api.Debug = b.Config.Debug
		}
	})

	return b.api, err
}

func (b *TelegramBot) Listen(ctx context.Context) (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	u.AllowedUpdates = []string{"message", "callback_query"}

	go func() {
		<-ctx.Done()
		b.api.StopReceivingUpdates()
	}()

	return b.api.GetUpdatesChan(u), nil
}
