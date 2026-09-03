package heartbeat

import (
	"context"
	"time"

	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"go.uber.org/zap"
)

// Service implements the heartbeat functionality.
type Service struct {
	config      Config
	monitorSvc  *monitor.Service
	messagesSvc *messages.Service
	bot         *telegram.Bot
	logger      *zap.Logger
}

// New creates a new heartbeat Service.
func New(
	config Config,
	monitorSvc *monitor.Service,
	messagesSvc *messages.Service,
	bot *telegram.Bot,
	logger *zap.Logger,
) *Service {
	return &Service{
		config:      config,
		monitorSvc:  monitorSvc,
		messagesSvc: messagesSvc,
		bot:         bot,
		logger:      logger,
	}
}

// Run starts the heartbeat service and blocks until context is canceled.
func (s *Service) Run(ctx context.Context) error {
	if !s.config.Enabled {
		s.logger.Info("heartbeat disabled")
		return nil
	}

	s.logger.Info("heartbeat started",
		zap.Duration("interval", s.config.Interval))

	ticker := time.NewTicker(s.config.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			s.sendHeartbeat()
		case <-ctx.Done():
			s.logger.Info("heartbeat stopped")
			return nil
		}
	}
}

func (s *Service) sendHeartbeat() {
	statuses := s.monitorSvc.GetCurrentStatuses()

	var online, offline int
	for _, st := range statuses {
		if st.State == monitor.ServiceStateOnline {
			online++
		} else {
			offline++
		}
	}

	msg, err := s.messagesSvc.Heartbeat(messages.HeartbeatContext{
		TotalServices:   len(statuses),
		OnlineServices:  online,
		OfflineServices: offline,
		CheckedAt:       time.Now(),
	})
	if err != nil {
		s.logger.Error("failed to render heartbeat message", zap.Error(err))
		return
	}

	chatID := s.config.ChatID
	if _, sendErr := s.bot.SendMessage(chatID, msg); sendErr != nil {
		s.logger.Error("failed to send heartbeat", zap.Error(sendErr))
	}
}
