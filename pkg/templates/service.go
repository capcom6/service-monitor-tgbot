package templates

import (
	"fmt"
	"strings"
	"sync"
	"text/template"
)

type Service struct {
	templates map[string]string
	escapeFn  func(string) string

	funcs template.FuncMap
	cache map[string]*template.Template
	mux   sync.RWMutex
}

func NewService(cfg Config) *Service {
	s := &Service{
		templates: cfg.Templates,
		escapeFn:  cfg.EscapeFn,

		funcs: nil,
		cache: make(map[string]*template.Template),
		mux:   sync.RWMutex{},
	}

	s.funcs = template.FuncMap{
		"escape": s.escape,
	}

	return s
}

func (s *Service) SetEscapeFn(fn func(string) string) {
	s.escapeFn = fn
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

	// Re-check cache in case another goroutine cached it while we waited for the lock
	if tmpl, ok := s.cache[name]; ok {
		return tmpl, nil
	}

	tmpl, err := template.New(name).Funcs(s.funcs).Parse(s.templates[name])
	if err != nil {
		return nil, fmt.Errorf("can't parse template %s: %w", name, err)
	}

	s.cache[name] = tmpl

	return tmpl, nil
}

func (s *Service) escape(text string) string {
	if s.escapeFn != nil {
		return s.escapeFn(text)
	}

	return text
}
