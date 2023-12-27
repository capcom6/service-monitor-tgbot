package telegram

import (
	"go.uber.org/fx"
	"go.uber.org/zap"
)

var Module = fx.Module(
	"telegram",
	fx.Decorate(func(log *zap.Logger) *zap.Logger {
		return log.Named("telegram")
	}),
	fx.Provide(NewTelegramBot),
)
