package messages

import (
	"github.com/capcom6/service-monitor-tgbot/pkg/templates"
	"go.uber.org/zap"
)

type Service struct {
	templatesSvc *templates.Service

	logger *zap.Logger
}

func NewService(cfg Config, logger *zap.Logger) *Service {
	if cfg.Templates == nil {
		cfg.Templates = make(map[string]string)
	}

	for k, v := range messageTemplates {
		if _, ok := cfg.Templates[k]; !ok {
			cfg.Templates[k] = v
		}
	}

	s := &Service{
		templatesSvc: templates.NewService(templates.Config{
			Templates: cfg.Templates,
			EscapeFn:  cfg.EscapeFn,
		}),

		logger: logger,
	}

	return s
}

func (s *Service) SetEscapeFn(fn func(string) string) {
	s.templatesSvc.SetEscapeFn(fn)
}

func (s *Service) Online(data OnlineContext) (string, error) {
	return s.templatesSvc.Render(TemplateOnline, data)
}

func (s *Service) Offline(data OfflineContext) (string, error) {
	return s.templatesSvc.Render(TemplateOffline, data)
}

func (s *Service) ServicesList(data ServicesListContext) (string, error) {
	return s.templatesSvc.Render(TemplateServicesList, data)
}
