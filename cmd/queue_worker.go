package main

import (
	"github.com/touhonoob/quikwallet/app"
	apiv1wallets "github.com/touhonoob/quikwallet/app/api/v1/wallets"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			app.NewLoggerMiddleware,
			app.NewDb,
			app.NewRedis,
			apiv1wallets.NewQueue,
			apiv1wallets.NewWalletsCache,
			apiv1wallets.NewWalletRepository,
		),
		fx.Invoke(
			func(lifecycle fx.Lifecycle, q apiv1wallets.IQueue) {
				if err := q.ConsumeNewWalletLogJob(); err != nil {
					panic(err)
				}
			},
		),
	).Run()
}