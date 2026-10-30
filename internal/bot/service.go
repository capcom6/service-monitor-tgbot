package bot

import (
	"context"
	"fmt"

	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/internal/sender"
	"go.uber.org/zap"
)

type Service struct {
	cfg Config

	senderSvc   *sender.Service
	monitorSvc  *monitor.Service
	messagesSvc *messages.Service

	logger *zap.Logger
}

func NewService(
	cfg Config,
	senderSvc *sender.Service,
	monitorSvc *monitor.Service,
	messagesSvc *messages.Service,
	logger *zap.Logger,
) *Service {
	return &Service{
		cfg: cfg,

		senderSvc:   senderSvc,
		monitorSvc:  monitorSvc,
		messagesSvc: messagesSvc,

		logger: logger,
	}
}

func (s *Service) Run(ctx context.Context) error {
	ch, err := s.monitorSvc.Monitor(ctx)
	if err != nil {
		return fmt.Errorf("failed to start monitoring services: %w", err)
	}

	for v := range ch {
		s.logger.Debug("probe", zap.String("name", v.Name), zap.String("state", string(v.State)), zap.Error(v.Error))

		duration := messages.FormatDurationSince(v.ChangedAt)

		var msg string
		if v.State == monitor.ServiceStateOffline {
			msg, err = s.messagesSvc.Offline(
				messages.OfflineContext{
					OnlineContext: messages.OnlineContext{
						Name:      v.Name,
						ChangedAt: v.ChangedAt,
						Duration:  duration,
					},
					Error: v.Error.Error(),
				},
			)
		} else {
			msg, err = s.messagesSvc.Online(messages.OnlineContext{
				Name:      v.Name,
				ChangedAt: v.ChangedAt,
				Duration:  duration,
			})
		}

		if err != nil {
			s.logger.Error("failed to render message template",
				zap.String("service_name", v.Name),
				zap.String("service_id", v.ID),
				zap.Error(err))
			continue
		}

		if sendErr := s.senderSvc.Send(ctx, s.cfg.ChatID, msg); sendErr != nil {
			s.logger.Error("failed to send notification",
				zap.String("service_name", v.Name),
				zap.String("service_id", v.ID),
				zap.Error(sendErr))
		}
	}

	return nil
}
