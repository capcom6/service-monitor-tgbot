package bot

import (
	"context"
	"fmt"

	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"go.uber.org/zap"
)

type Service struct {
	cfg Config

	bot      *telegram.Bot
	monitor  *monitor.MonitorModule
	messages *messages.Service

	logger *zap.Logger
}

func NewService(cfg Config, bot *telegram.Bot, monitor *monitor.MonitorModule, messages *messages.Service, logger *zap.Logger) *Service {
	return &Service{
		cfg: cfg,

		bot:      bot,
		monitor:  monitor,
		messages: messages,

		logger: logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	ch, err := s.monitor.Monitor(ctx)
	if err != nil {
		return fmt.Errorf("failed to monitor services: %w", err)
	}

	for v := range ch {
		s.logger.Debug("probe", zap.String("name", v.Name), zap.String("state", string(v.State)), zap.Error(v.Error))

		msg := ""
		if v.State == monitor.ServiceStateOffline {
			context := OfflineContext{
				OnlineContext: OnlineContext{
					Name: s.bot.EscapeText(v.Name),
				},
				Error: s.bot.EscapeText(v.Error.Error()),
			}
			msg, err = s.messages.Render(TemplateOffline, context)
		} else {
			context := OnlineContext{
				Name: s.bot.EscapeText(v.Name),
			}
			msg, err = s.messages.Render(TemplateOnline, context)
		}

		if err != nil {
			s.logger.Error("can't render template", zap.Error(err))
			continue
		}

		if _, err := s.bot.SendMessage(s.cfg.ChatID, msg); err != nil {
			s.logger.Error("can't send message", zap.Error(err))
		}
	}

	return nil
}

func (s *Service) HandleStatusCommand(ctx context.Context, cmd telegram.Command) {
	services := s.monitor.GetCurrentStatuses()

	if len(services) == 0 {
		if _, err := s.bot.SendMessage(cmd.From, "No services configured."); err != nil {
			s.logger.Error("can't send message", zap.Error(err))
		}
		return
	}

	context := make(ServicesListContext, len(services))

	for i, service := range services {
		context[i] = ServiceState{
			Name:  service.Name,
			State: string(service.State),
			Error: service.Error,
		}
	}

	msg, err := s.messages.Render(TemplateServicesList, context)
	if err != nil {
		s.logger.Error("can't render template", zap.Error(err))
		return
	}

	if _, err := s.bot.SendMessage(cmd.From, msg); err != nil {
		s.logger.Error("can't send message", zap.Error(err))
	}
}
