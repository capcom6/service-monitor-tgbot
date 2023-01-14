package monitor

import "context"

type Pinger interface {
	Ping(context.Context) error
}
