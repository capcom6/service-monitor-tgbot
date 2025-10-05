package messages

import (
	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"messages",
		fxutil.WithNamedLogger("messages"),
		fx.Provide(NewService),
	)
}
