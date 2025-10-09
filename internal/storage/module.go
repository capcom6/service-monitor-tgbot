package storage

import (
	"os"

	"github.com/capcom6/go-infra-fx/fxutil"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"storage",
		fxutil.WithNamedLogger("storage"),
		fx.Provide(func() Storage {
			return NewService(
				&yamlStorage{
					Path: os.Getenv("CONFIG_PATH"),
				},
			)
		}),
	)
}
