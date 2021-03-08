package main

import (
	"fmt"
	"github.com/shopspring/decimal"
	"github.com/touhonoob/quikwallet/app"
	apiv1wallets "github.com/touhonoob/quikwallet/app/api/v1/wallets"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/gavv/httpexpect/v2"
)

type WalletsApiClient struct {
	t *testing.T
}

func (client *WalletsApiClient) newHttpExpect() *httpexpect.Expect {
	return httpexpect.New(
		client.t,
		fmt.Sprintf("http://quikwallet:%s/api/v1", os.Getenv("QUIKWALLET_PORT")),
	)
}

func (client *WalletsApiClient) getBalance() *httpexpect.Response {
	return client.newHttpExpect().GET(
		fmt.Sprintf(
			"/wallets/%s/balance",
			os.Getenv("QUIKWALLET_PREPOPULATED_WALLET_UUID"),
		),
	).Expect()
}

func (client *WalletsApiClient) postCredit(credit string) *httpexpect.Response {
	return client.newHttpExpect().POST(
		fmt.Sprintf(
			"/wallets/%s/credit",
			os.Getenv("QUIKWALLET_PREPOPULATED_WALLET_UUID"),
		),
	).WithJSON(
		map[string]string{
			"credit": credit,
		},
	).Expect()
}

func (client *WalletsApiClient) postDebit(debit string) *httpexpect.Response {
	return client.newHttpExpect().POST(
		fmt.Sprintf(
			"/wallets/%s/debit",
			os.Getenv("QUIKWALLET_PREPOPULATED_WALLET_UUID"),
		),
	).WithJSON(
		map[string]string{
			"debit": debit,
		},
	).Expect()
}

func (client *WalletsApiClient) getWalletLog(uuid string) *httpexpect.Response {
	return client.newHttpExpect().GET(
		fmt.Sprintf(
			"/wallets/%s/logs/%s",
			os.Getenv("QUIKWALLET_PREPOPULATED_WALLET_UUID"), uuid,
		),
	).Expect()
}

func setUp() {
	db := app.NewDb()
	if err := db.Exec("TRUNCATE wallet_logs;").Error; err != nil {
		panic(err)
	}
}

func TestGetBalance(t *testing.T) {
	setUp()

	client := &WalletsApiClient{t: t}
	client.getBalance().Status(http.StatusOK).JSON().Object().ContainsKey("balance")
}

func TestCreditAndDebit(t *testing.T) {
	setUp()

	client := &WalletsApiClient{t: t}
	const credit = "12.34"
	originalBalance := client.getBalance().JSON().Object().Value("balance").String().Raw()
	walletLogUuid := client.postCredit(credit).Status(http.StatusAccepted).JSON().Object().ContainsKey("wallet_log_uuid").Value("wallet_log_uuid").String().Raw()

	time.Sleep(time.Duration(500) * time.Millisecond)

	client.getWalletLog(walletLogUuid).Status(http.StatusOK).JSON().Object().ValueEqual(
		"status", 1,
	)
	client.getBalance().Status(http.StatusOK).JSON().Object().ValueEqual(
		"balance",
		decimal.RequireFromString(originalBalance).Add(decimal.RequireFromString(credit)).String(),
	)

	const debit = "12.34"
	originalBalance = client.getBalance().JSON().Object().Value("balance").String().Raw()
	walletLogUuid = client.postDebit(debit).Status(http.StatusAccepted).JSON().Object().ContainsKey("wallet_log_uuid").Value("wallet_log_uuid").String().Raw()

	time.Sleep(time.Duration(500) * time.Millisecond)

	client.getWalletLog(walletLogUuid).Status(http.StatusOK).JSON().Object().ValueEqual(
		"status", apiv1wallets.Accepted,
	)
	client.getBalance().Status(http.StatusOK).JSON().Object().ValueEqual(
		"balance",
		decimal.RequireFromString(originalBalance).Sub(decimal.RequireFromString(debit)).String(),
	)
}

func TestDebitHuge(t *testing.T) {
	setUp()

	client := &WalletsApiClient{t: t}
	const debit = "999999"
	originalBalance := client.getBalance().JSON().Object().Value("balance").String().Raw()
	walletLogUuid := client.postDebit(debit).Status(http.StatusAccepted).JSON().Object().ContainsKey("wallet_log_uuid").Value("wallet_log_uuid").String().Raw()

	time.Sleep(time.Duration(500) * time.Millisecond)

	client.getWalletLog(walletLogUuid).Status(http.StatusOK).JSON().Object().ValueEqual(
		"status", apiv1wallets.Rejected,
	)
	client.getBalance().Status(http.StatusOK).JSON().Object().ValueEqual(
		"balance",
		originalBalance,
	)
}
