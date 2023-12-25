package storage

import (
	"os"

	"go.uber.org/fx"
)

var Module = fx.Module(
	"storage",
	fx.Provide(func() Storage {
		return NewStorageService(
			&yamlStorage{
				Path: os.Getenv("CONFIG_PATH"),
			},
		)
	}),
)
