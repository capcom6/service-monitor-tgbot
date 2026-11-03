package bot

import (
	"github.com/capcom6/service-monitor-tgbot/internal/bot/handler"
	"github.com/capcom6/service-monitor-tgbot/internal/bot/handlers/status"
	"github.com/go-core-fx/fxutil"
	"github.com/go-core-fx/logger"
	"github.com/go-core-fx/telegofx"
	"go.uber.org/fx"
)

func Module() fx.Option {
	return fx.Module(
		"bot",
		logger.WithNamedLogger("bot"),
		fx.Provide(NewService),
		fx.Invoke(fxutil.RegisterRunnable[*Service]()),
		// One subpackage per handler: each constructor is annotated into the
		// "handlers" value group and returns handler.Handler, the type dig
		// keys the group on. Registration order on the router follows provide
		// order; telegohandler dispatch is first-match-wins.
		fx.Provide(
			fx.Annotate(status.New, fx.ResultTags(`group:"handlers"`)),
			fx.Private,
		),
		fx.Invoke(fx.Annotate(
			func(handlers []handler.Handler, router *telegofx.Router) {
				for _, h := range handlers {
					h.Register(router)
				}
			},
			fx.ParamTags(`group:"handlers"`),
		)),
	)
}
