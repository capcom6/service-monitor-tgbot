package events

import "encoding/json"

type EventName string
type ServiceState string

const (
	ServiceStateOnline  = "online"
	ServiceStateOffline = "offline"

	EventNameServiceStateChanged = "service_state_changed"
)

type Event[T any] struct {
	Name    EventName `json:"name"`
	Payload T         `json:"data"`
}

func (e Event[T]) Encode() (string, error) {
	data, err := json.Marshal(e)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

func (e *Event[T]) Decode(input string) error {
	return json.Unmarshal([]byte(input), e)
}

type ServiceStateChanged struct {
	Id    string `json:"id"`
	Name  string `json:"name"`
	State string `json:"state"`
	Error string `json:"error,omitempty"`
}

func NewServiceStateChangedEvent(id, name, state string, err error) Event[ServiceStateChanged] {
	errtext := ""
	if err != nil {
		errtext = err.Error()
	}

	return Event[ServiceStateChanged]{
		Name: EventNameServiceStateChanged,
		Payload: ServiceStateChanged{
			Id:    id,
			Name:  name,
			State: state,
			Error: errtext,
		},
	}
}
