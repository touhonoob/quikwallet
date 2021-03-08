package apiv1wallets

import (
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type IApiV1Wallets interface {
	Router() gin.IRouter
}

type IWalletsController interface {
	PostCredit(c *gin.Context)
	PostDebit(c *gin.Context)
	GetBalance(c *gin.Context)
	GetWalletLog(c *gin.Context)
}

type IWalletsRepository interface {
	CreateCreditLog(walletUUID uuid.UUID, credit decimal.Decimal) (*WalletLog, error)
	CreateDebitLog(walletUUID uuid.UUID, debit decimal.Decimal) (*WalletLog, error)
	ProcessWalletLogs(walletUUID uuid.UUID) error
	GetWallet(walletUUID uuid.UUID) (*Wallet, error)
	GetBalance(walletUUID uuid.UUID) (decimal.Decimal, error)
	GetWalletLog(walletLogUUID uuid.UUID) (*WalletLog, error)
}

type IWalletsCache interface {
	GetWalletBalance(walletUUID uuid.UUID) (decimal.Decimal, error)
	InvalidateWalletBalance(walletUUID uuid.UUID) error
}

type IQueue interface {
	PublishNewWalletLogJob(walletUUID uuid.UUID) error
	ConsumeNewWalletLogJob() error
}