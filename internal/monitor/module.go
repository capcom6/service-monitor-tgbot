package monitor

import (
	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"monitor",
		fxutil.WithNamedLogger("monitor"),
		fx.Provide(NewService),
	)
}
