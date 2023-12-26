package infrastructure

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"infrastructure",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("infrastructure")
	}),
	fx.Provide(NewTelegramBot),
)
