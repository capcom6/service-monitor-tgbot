package server

import (
	"github.com/capcom6/service-monitor-tgbot/internal/server/docs"
	"github.com/go-core-fx/fiberfx"
	"github.com/go-core-fx/fiberfx/handler"
	"github.com/go-core-fx/fiberfx/health"
	"github.com/go-core-fx/fiberfx/openapi"
	"github.com/go-core-fx/fiberfx/validation"
	"github.com/go-core-fx/logger"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

func Module() fx.Option {
	return fx.Module(
		"server",
		logger.WithNamedLogger("server"),

		fx.Provide(func(log *zap.Logger) fiberfx.Options {
			opts := fiberfx.Options{}
			opts.WithErrorHandler(fiberfx.NewJSONErrorHandler(log))
			opts.WithMetrics()
			return opts
		}),
		fx.Supply(docs.SwaggerInfo),

		fx.Provide(
			health.NewHandler,
			openapi.NewHandler,
			fx.Private,
		),

		// fx.Provide(
		// 	fx.Annotate(example.New, fx.ResultTags(`group:"handlers"`)),
		// 	fx.Private,
		// ),

		fx.Invoke(
			fx.Annotate(
				func(handlers []handler.Handler, healthHandler *health.Handler, openapiHandler *openapi.Handler, app *fiber.App) {
					// Health endpoint
					healthHandler.Register(app)

					// Version 1 API group
					v1 := app.Group("/api/v1")
					openapiHandler.Register(v1.Group("/docs"))

					v1.Use(validation.Middleware)

					for _, h := range handlers {
						h.Register(v1)
					}
				},
				fx.ParamTags(`group:"handlers"`),
			),
		),
	)
}
