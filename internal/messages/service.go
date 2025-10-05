package messages

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	"go.uber.org/zap"
)

type Service struct {
	templates map[string]string

	cache map[string]*template.Template
	mux   sync.RWMutex

	logger *zap.Logger
}

func NewService(cfg Config, logger *zap.Logger) *Service {
	return &Service{
		templates: cfg.Templates,

		cache: make(map[string]*template.Template),
		mux:   sync.RWMutex{},

		logger: logger,
	}
}

func (s *Service) Render(name string, data any) (string, error) {
	tmpl, err := s.prepare(name)
	if err != nil {
		return "", fmt.Errorf("can't prepare template: %w", err)
	}

	builder := strings.Builder{}
	if err := tmpl.Execute(&builder, data); err != nil {
		return "", fmt.Errorf("can't execute template: %w", err)
	}

	return builder.String(), nil
}

func (s *Service) prepare(name string) (*template.Template, error) {
	if _, ok := s.templates[name]; !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	s.mux.RLock()
	if tmpl, ok := s.cache[name]; ok {
		s.mux.RUnlock()
		return tmpl, nil
	}
	s.mux.RUnlock()

	s.mux.Lock()
	defer s.mux.Unlock()

	tmpl, err := template.New(name).Parse(s.templates[name])
	if err != nil {
		return nil, fmt.Errorf("can't parse template %s: %w", name, err)
	}

	s.cache[name] = tmpl

	return tmpl, nil
}
