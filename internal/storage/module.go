package storage

import (
	"os"

	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"storage",
		logger.WithNamedLogger("storage"),
		fx.Provide(func() Storage {
			return NewService(
				&yamlStorage{
					Path: os.Getenv("CONFIG_PATH"),
				},
			)
		}),
	)
}
