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

	for _, v := range cfg.Services {
		if err := monitor.NewMonitorService(v, tgapi).Start(ctx); err != nil {
			errorLog.Printf("Can't monitor service %s: %s\n", v.Name, err.Error())
		}
	}

	log.Println("Started")

	<-ctx.Done()

	log.Println("Done")

	return nil
}
