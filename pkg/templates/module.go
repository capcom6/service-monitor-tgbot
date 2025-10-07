package templates

import (
	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"templates",
		fxutil.WithNamedLogger("templates"),
		fx.Provide(NewService),
	)
}
