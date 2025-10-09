package messages

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
	Error string
}

func NewServiceState(name, state string, err error) ServiceState {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	return ServiceState{
		Name:  name,
		State: state,
		Error: errStr,
	}
}

type ServicesListContext []ServiceState
