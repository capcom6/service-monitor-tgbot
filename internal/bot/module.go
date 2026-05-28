package bot

import (
	"github.com/capcom6/service-monitor-tgbot/internal/messages"
	"github.com/capcom6/service-monitor-tgbot/pkg/telegram"
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(NewService),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
		fx.Invoke(func(
			s *Service,
			tg *telegram.Bot,
			messages *messages.Service,
		) {
			messages.SetEscapeFn(tg.EscapeText)
			tg.AddHandler("status", s.HandleStatusCommand)
		}),
	)
}
