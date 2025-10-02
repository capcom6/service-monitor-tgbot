package bot

import (
	"fmt"
	"strings"
	"sync"
	"text/template"

	"github.com/capcom6/service-monitor-tgbot/internal/config"
	"github.com/capcom6/service-monitor-tgbot/internal/monitor"
)

type TemplateName string

const (
	TemplateOnline  TemplateName = "online"
	TemplateOffline TemplateName = "offline"
)

type Messages struct {
	Templates       config.TelegramMessages
	cachedTemplates map[TemplateName]*template.Template
	mux             sync.Mutex
}

type OnlineContext struct {
	Name string
}

type OfflineContext struct {
	OnlineContext
	Error string
}

func NewMessages(templates config.TelegramMessages) *Messages {
	return &Messages{
		Templates:       templates,
		cachedTemplates: make(map[TemplateName]*template.Template),
	}
}

func (m *Messages) prepare(name TemplateName) (prepared *template.Template, err error) {
	m.mux.Lock()
	defer m.mux.Unlock()

	if prepared = m.cachedTemplates[name]; prepared != nil {
		return prepared, nil
	}

	var tmplString string
	var ok bool
	if tmplString, ok = m.Templates[string(name)]; !ok {
		return nil, fmt.Errorf("template %s not found", name)
	}

	if m.cachedTemplates[name], err = template.New(string(name)).Parse(tmplString); err != nil {
		return nil, fmt.Errorf("can't parse template %s: %w", name, err)
	}

	return m.cachedTemplates[name], nil
}

func (m *Messages) Render(name TemplateName, context any) (string, error) {
	tmpl, err := m.prepare(name)
	if err != nil {
		return "", err
	}

	builder := strings.Builder{}
	if err := tmpl.Execute(&builder, context); err != nil {
		return "", err
	}
	return builder.String(), nil
}

func (m *Messages) RenderStatus(services []monitor.ServiceStatus) string {
	if len(services) == 0 {
		return "No services configured."
	}

	var builder strings.Builder

	for i, service := range services {
		if i > 0 {
			builder.WriteString("\n")
		}

		stateEmoji := "❌"
		if service.State == monitor.ServiceOnline {
			stateEmoji = "✅"
		}

		stateText := "OFFLINE"
		if service.State == monitor.ServiceOnline {
			stateText = "ONLINE"
		}

		// Format: "Service: [state] (last checked: [time])"
		// Example: "API Server: ONLINE ✅ (2m ago)"
		builder.WriteString(fmt.Sprintf("%s: %s %s", service.Name, stateText, stateEmoji))

		// Add last checked time if available
		if service.Error != nil {
			builder.WriteString(fmt.Sprintf(" (%s)", service.Error.Error()))
		} else {
			builder.WriteString(" (recently)")
		}
	}

	return builder.String()
}
