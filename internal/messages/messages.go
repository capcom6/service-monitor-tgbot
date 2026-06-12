package messages

import (
	"time"
)

const (
	TemplateOnline       string = "online"
	TemplateOffline      string = "offline"
	TemplateServicesList string = "services_list"
)

type OnlineContext struct {
	Name      string
	ChangedAt time.Time
	Duration  string
}

type OfflineContext struct {
	OnlineContext

	Error string
}

type ServiceState struct {
	Name      string
	State     string
	Error     string
	ChangedAt time.Time
	Duration  string
}

func NewServiceState(name, state string, err error, changedAt time.Time) ServiceState {
	errStr := ""
	if err != nil {
		errStr = err.Error()
	}

	duration := FormatDurationSince(changedAt)

	return ServiceState{
		Name:      name,
		State:     state,
		Error:     errStr,
		ChangedAt: changedAt,
		Duration:  duration,
	}
}

type ServicesListContext []ServiceState

func FormatDurationSince(t time.Time) string {
	if t.IsZero() {
		return "unknown"
	}
	return time.Since(t).Truncate(time.Second).String()
}
