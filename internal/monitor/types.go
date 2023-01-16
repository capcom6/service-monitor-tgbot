package monitor

import "context"

type Probeer interface {
	Probe(context.Context) error
}
