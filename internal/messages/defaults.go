//nolint:lll //one-line messages templates
package messages

func defaultTemplates() map[string]string {
	return map[string]string{
		TemplateOnline:       "✅ {{.Name | escape}} is *online*",
		TemplateOffline:      "❌ {{.Name | escape}} is *offline*: {{.Error | escape}}",
		TemplateServicesList: "📋 Services:\n\n{{range .}}\\- {{.Name | escape}}: {{if eq .State \"online\"}}✅{{else}}❌{{with .Error}} \\({{. | escape}}\\){{end}}{{end}}\n{{end}}",
	}
}
