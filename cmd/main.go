package main

import (
	"github.com/touhonoob/quikwallet/app"
	"github.com/touhonoob/quikwallet/app/api/v1"
	"github.com/touhonoob/quikwallet/app/api/v1/wallets"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			app.NewLoggerMiddleware,
			app.NewDb,
			app.NewRedis,
			app.NewApp,
			apiv1.NewApiV1,
			apiv1wallets.NewQueue,
			apiv1wallets.NewWalletsCache,
			apiv1wallets.NewApiV1Wallets,
			apiv1wallets.NewWalletRepository,
			apiv1wallets.NewWalletsController,
		),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, app app.IQuikWalletApp, apiV1 apiv1.IApiV1, apiV1Wallets apiv1wallets.IApiV1Wallets) {
				if err := app.Run(); err != nil {
					panic(err)
				}
			},
		),
	).Run()
}
