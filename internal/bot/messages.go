package bot

const (
	TemplateOnline  string = "online"
	TemplateOffline string = "offline"
)

type OnlineContext struct {
	Name string
}

type OfflineContext struct {
	OnlineContext
	Error string
}

// func (m *Messages) RenderStatus(services []monitor.ServiceStatus) string {
// 	if len(services) == 0 {
// 		return "No services configured."
// 	}

// 	var builder strings.Builder

// 	for i, service := range services {
// 		if i > 0 {
// 			builder.WriteString("\n")
// 		}

// 		stateEmoji := "❌"
// 		if service.State == monitor.ServiceOnline {
// 			stateEmoji = "✅"
// 		}

// 		stateText := "OFFLINE"
// 		if service.State == monitor.ServiceOnline {
// 			stateText = "ONLINE"
// 		}

// 		// Format: "Service: [state] (last checked: [time])"
// 		// Example: "API Server: ONLINE ✅ (2m ago)"
// 		builder.WriteString(fmt.Sprintf("%s: %s %s", service.Name, stateText, stateEmoji))

// 		// Add last checked time if available
// 		if service.Error != nil {
// 			builder.WriteString(fmt.Sprintf(" (%s)", service.Error.Error()))
// 		} else {
// 			builder.WriteString(" (recently)")
// 		}
// 	}

// 	return builder.String()
// }
