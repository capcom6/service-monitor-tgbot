package monitor

import (
	"context"
	"time"

	"github.com/capcom6/tgbot-service-monitor/internal/config"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type MonitorService struct {
	Bot     *tgbotapi.BotAPI
	Service config.Service

	p Pinger
}

func NewMonitorService(service config.Service, bot *tgbotapi.BotAPI) *MonitorService {
	svc := MonitorService{
		Bot:     bot,
		Service: service,
	}

	if !service.HTTPGet.IsEmpty() {
		svc.p = NewHttpPinger(service.HTTPGet)
	}

	return &svc
}

func (s *MonitorService) Start(ctx context.Context) (err error) {
	if s.Service, err = s.Service.ApplyDefaultsAndValidate(); err != nil {
		return
	}

	go func() {
		log.Println("Init of", s.Service.Name)

		if s.Service.InitialDelaySeconds > 0 {
			time.Sleep(time.Duration(s.Service.InitialDelaySeconds) * time.Second)
		}

		log.Println("Start of", s.Service.Name)
		ticker := time.NewTicker(time.Duration(s.Service.PeriodSeconds) * time.Second)
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:
				c, cancel := context.WithTimeout(ctx, time.Duration(s.Service.TimeoutSeconds)*time.Second)
				err = s.p.Ping(c)
				if err != nil {
					errorLog.Println("Error of", s.Service.Name, ":", err.Error())
				}
				cancel()

				log.Println("Tick for", s.Service.Name)
			case <-ctx.Done():
				log.Println("Stop of", s.Service.Name)
				return
			}
		}
	}()

	return
}
