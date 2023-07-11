package bot

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/infrastructure"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() error {
	rand.Seed(time.Now().UnixNano())

	cfg := config.GetConfig()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

	messages := NewMessages(cfg.Telegram.Messages)
	tg := infrastructure.NewTelegramBot(cfg.Telegram)
	tgapi, err := tg.Api()
	if err != nil {
		return fmt.Errorf("can't init Telegram Api client: %w", err)
	}

	module := monitor.NewMonitorModule(cfg.Services)
	ch, err := module.Monitor(ctx)
	if err != nil {
		return err
	}

	log.Println("Started")

	for v := range ch {
		log.Printf("%+v\n", v)

		msg := tgbotapi.NewMessage(cfg.Telegram.ChatID, "")
		msg.ParseMode = tgbotapi.ModeMarkdownV2
		if v.State == monitor.ServiceOffline {
			// msg.Text = "❌ " + tgbotapi.EscapeText(msg.ParseMode, v.Name) + " is *offline*: " + tgbotapi.EscapeText(msg.ParseMode, v.Error.Error())
			context := OfflineContext{
				OnlineContext: OnlineContext{
					Name: tgbotapi.EscapeText(msg.ParseMode, v.Name),
				},
				Error: tgbotapi.EscapeText(msg.ParseMode, v.Error.Error()),
			}
			msg.Text, err = messages.Render(TemplateOffline, context)
		} else {
			// msg.Text = "✅ " + tgbotapi.EscapeText(msg.ParseMode, v.Name) + " is *online*"
			context := OnlineContext{
				Name: tgbotapi.EscapeText(msg.ParseMode, v.Name),
			}
			msg.Text, err = messages.Render(TemplateOnline, context)
		}

		if err != nil {
			errorLog.Println(err)
			continue
		}

		if _, err := tgapi.Send(msg); err != nil {
			errorLog.Println(err)
		}
	}

	<-ctx.Done()

	log.Println("Done")

	return nil
}
