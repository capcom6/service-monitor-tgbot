package bot

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"time"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
	"github.com/capcom6/tgbot-service-monitor/internal/infrastructure"
	"github.com/capcom6/tgbot-service-monitor/internal/monitor"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

func Run() error {
	rand.Seed(time.Now().UnixNano())

	cfg := config.GetConfig()
	ctx, _ := signal.NotifyContext(context.Background(), os.Interrupt)

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
		msg := tgbotapi.NewMessage(cfg.Telegram.ChatID, fmt.Sprintf("%+v", v))
		if _, err := tgapi.Send(msg); err != nil {
			errorLog.Println(err)
		}
		// tgapi.Se
	}

	// for _, v := range cfg.Services {
	// 	if err := monitor.NewMonitorService(v, tgapi).Start(ctx); err != nil {
	// 		errorLog.Printf("Can't monitor service %s: %s\n", v.Name, err.Error())
	// 	}
	// }

	<-ctx.Done()

	log.Println("Done")

	return nil
}
