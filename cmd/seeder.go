package main

import (
	"fmt"
	"github.com/touhonoob/quikwallet/app"
	apiv1wallets "github.com/touhonoob/quikwallet/app/api/v1/wallets"
	"os"
	"time"
	"github.com/rs/zerolog/log"
)

func main () {
	var prepopulated int64
	db := app.NewDb()
	if err := db.Model(&apiv1wallets.Wallet{}).Where(
		"uuid = ?", PrepopulatedWalletUuid(),
	).Count(&prepopulated).Error; err != nil {
		panic(err)
	} else if prepopulated != 0 {
		log.Printf("wallet %s is already created\n", PrepopulatedWalletUuid())
	} else if err := db.Model(&apiv1wallets.Wallet{}).Create(&apiv1wallets.Wallet{
		Uuid: PrepopulatedWalletUuid(),
		CreatedAt: time.Now(),
	}).Error; err != nil {
		panic(err)
	} else {
		fmt.Printf("created wallet %s\n", PrepopulatedWalletUuid())
	}
}

func PrepopulatedWalletUuid() string {
	return os.Getenv("QUIKWALLET_PREPOPULATED_WALLET_UUID")
}