package monitor

import "context"

const (
	ServiceOnline  ServiceState = "online"
	ServiceOffline ServiceState = "offline"
)

type ServiceState string

type Probeer interface {
	Probe(context.Context) error
}

type ServiceStatus struct {
	Id    string
	Name  string
	State ServiceState
	Error error
}

type ServiceProbe struct {
	Id    string
	Name  string
	Error error
}

type ProbesChannel chan ServiceProbe
type UpdatesChannel chan ServiceStatus
