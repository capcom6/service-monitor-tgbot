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
	monitor  *monitor.Service
	messages *messages.Service

	logger *zap.Logger
}

func NewService(
	cfg Config,
	bot *telegram.Bot,
	monitor *monitor.Service,
	messages *messages.Service,
	logger *zap.Logger,
) *Service {
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

		var msg string
		if v.State == monitor.ServiceStateOffline {
			msg, err = s.messages.Offline(
				messages.OfflineContext{
					OnlineContext: messages.OnlineContext{
						Name: v.Name,
					},
					Error: v.Error.Error(),
				},
			)
		} else {
			msg, err = s.messages.Online(messages.OnlineContext{
				Name: v.Name,
			})
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

func (s *Service) HandleStatusCommand(_ context.Context, cmd telegram.Command) {
	services := s.monitor.GetCurrentStatuses()

	if len(services) == 0 {
		if _, err := s.bot.SendMessage(cmd.Chat, "No services configured."); err != nil {
			s.logger.Error("can't send message", zap.Error(err))
		}
		return
	}

	context := make(messages.ServicesListContext, len(services))

	for i, service := range services {
		context[i] = messages.NewServiceState(service.Name, string(service.State), service.Error)
	}

	msg, err := s.messages.ServicesList(context)
	if err != nil {
		s.logger.Error("can't render template", zap.Error(err))
		return
	}

	if _, err := s.bot.SendMessage(cmd.Chat, msg); err != nil {
		s.logger.Error("can't send message", zap.Error(err))
	}
}
