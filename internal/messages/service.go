package messages

import (
	"fmt"

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

	for k, v := range defaultTemplates() {
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
	return s.render(TemplateOnline, data)
}

func (s *Service) Offline(data OfflineContext) (string, error) {
	return s.render(TemplateOffline, data)
}

func (s *Service) ServicesList(data ServicesListContext) (string, error) {
	return s.render(TemplateServicesList, data)
}

func (s *Service) render(name string, data any) (string, error) {
	res, err := s.templatesSvc.Render(name, data)
	if err != nil {
		return "", fmt.Errorf("can't render template %s: %w", name, err)
	}

	return res, nil
}
