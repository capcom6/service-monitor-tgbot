package botx

import (
	"fmt"

	"go.uber.org/fx"
)

func Run() {
	module := fx.Module(
		"bot",
		fx.Invoke(func() {
			fmt.Println("Bot")
		}),
	)

	fx.New(module).Run()
}
