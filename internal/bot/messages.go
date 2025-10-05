package bot

const (
	TemplateOnline       string = "online"
	TemplateOffline      string = "offline"
	TemplateServicesList string = "services_list"
)

type OnlineContext struct {
	Name string
}

type OfflineContext struct {
	OnlineContext
	Error string
}

type ServiceState struct {
	Name  string
	State string
	Error error
}

type ServicesListContext []ServiceState
