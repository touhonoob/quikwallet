package main

import (
	"github.com/touhonoob/quikwallet/app"
	apiv1wallets "github.com/touhonoob/quikwallet/app/api/v1/wallets"
)

func main() {
	db := app.NewDb()
	if err := db.AutoMigrate(&apiv1wallets.Wallet{}, &apiv1wallets.WalletLog{}); err != nil {
		panic(err)
	}
}
