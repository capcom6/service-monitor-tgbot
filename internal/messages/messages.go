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

type serviceState struct {
	Name  string
	State string
	Error string
}

func NewServiceState(name, state string, err error) serviceState {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	return serviceState{
		Name:  name,
		State: string(state),
		Error: errStr,
	}
}

type ServicesListContext []serviceState
