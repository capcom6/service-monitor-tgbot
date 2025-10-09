package monitor

import (
	"context"
	"time"
)

type ServiceState string

const (
	ServiceStateUnknown ServiceState = "unknown"
	ServiceStateOnline  ServiceState = "online"
	ServiceStateOffline ServiceState = "offline"
)

type Probeer interface {
	Probe(context.Context) error
}

type ServiceStatus struct {
	ID    string
	Name  string
	State ServiceState
	Error error
}

type state struct {
	Probes    int
	Online    bool
	Error     error
	Timestamp time.Time
}

func (s state) State() ServiceState {
	if s.Timestamp.IsZero() {
		return ServiceStateUnknown
	}

	if s.Online {
		return ServiceStateOnline
	}

	return ServiceStateOffline
}

type ProbesChannel chan error
type UpdatesChannel chan ServiceStatus
