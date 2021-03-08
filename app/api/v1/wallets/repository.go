package apiv1wallets

import (
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
	"time"
)

type Wallet struct {
	Uuid       string      `gorm:"primaryKey; type:varchar(36) NOT NULL;"`
	WalletLogs []WalletLog `gorm:"references:Uuid;foreignKey:WalletUuid;constraint:OnUpdate:RESTRICT,OnDelete:RESTRICT;"`
	CreatedAt  time.Time
}

type WalletLog struct {
	Uuid       string          `gorm:"primaryKey; type:varchar(36) NOT NULL;" json:"uuid"`
	WalletUuid string          `gorm:"type:varchar(36) NOT NULL;" json:"wallet_uuid"`
	Log        int64           `gorm:"type:bigint(20) NOT NULL;" json:"log"`
	Status     WalletLogStatus `gorm:"type:tinyint(1) NOT NULL;default:0;" json:"status"`
	CreatedAt  time.Time       `json:"created_at"`
}

type WalletLogStatus uint8

const (
	ToBeProcessed WalletLogStatus = iota
	Accepted                      = iota
	Rejected                      = iota
)

type WalletsRepository struct {
	db *gorm.DB
}

func (repo *WalletsRepository) CreateCreditLog(walletUUID uuid.UUID, credit decimal.Decimal) (*WalletLog, error) {
	creditLog := &WalletLog{
		Uuid:       uuid.NewString(),
		WalletUuid: walletUUID.String(),
		Log:        credit.Shift(2).IntPart(),
		CreatedAt:  time.Now(),
	}
	return creditLog, repo.db.Create(creditLog).Error
}

func (repo *WalletsRepository) CreateDebitLog(walletUUID uuid.UUID, debit decimal.Decimal) (*WalletLog, error) {
	debitLog := &WalletLog{
		Uuid:       uuid.NewString(),
		WalletUuid: walletUUID.String(),
		Log:        debit.Shift(2).Neg().IntPart(),
		CreatedAt:  time.Now(),
	}
	return debitLog, repo.db.Create(debitLog).Error
}

func (repo *WalletsRepository) ProcessWalletLogs(walletUUID uuid.UUID) error {
	var logs []WalletLog
	if err := repo.db.Model(&WalletLog{}).
		Where("wallet_uuid = ?", walletUUID.String()).
		Order("created_at ASC").
		Find(&logs).Error; err != nil {
		return err
	} else {
		for _, walletLog := range logs {
			var balance decimal.Decimal
			if _balance, err := repo.GetBalance(walletUUID); err != nil {
				return err
			} else {
				balance = _balance
			}

			newBalance := balance.Add(decimal.NewFromInt(walletLog.Log)).Floor()
			log.Info().Msgf("%s + %s = %s (%b)", balance.String(), decimal.NewFromInt(walletLog.Log).String(), newBalance.String(), newBalance.IsNegative())
			if newBalance.IsNegative() {
				if err := repo.db.Model(&walletLog).Update(
					"status", Rejected,
				).Error; err != nil {
					return err
				}
			} else if err := repo.db.Model(&walletLog).Update(
				"status", Accepted,
			).Error; err != nil {
				return err
			} else {
				continue
			}
		}
		return nil
	}
}

func (repo *WalletsRepository) GetWallet(walletUUID uuid.UUID) (*Wallet, error) {
	var wallet Wallet
	if err := repo.db.Model(&Wallet{}).Where(
		"uuid = ?", walletUUID.String(),
	).First(&wallet).Error; err != nil {
		return nil, err
	} else {
		return &wallet, nil
	}
}

func (repo *WalletsRepository) GetBalance(walletUUID uuid.UUID) (decimal.Decimal, error) {
	var sum int64
	if err := repo.db.Model(&WalletLog{}).
		Select("IFNULL(SUM(log),0)").
		Where("wallet_uuid = ?", walletUUID.String()).
		Where("status = ?", Accepted).
		Scan(&sum).Error; err != nil {
		return decimal.Zero, err
	} else {
		return decimal.NewFromInt(sum), nil
	}
}

func (repo *WalletsRepository) GetWalletLog(walletLogUUID uuid.UUID) (*WalletLog, error) {
	var walletLog WalletLog
	if err := repo.db.Model(&WalletLog{}).Where(
		"uuid = ?", walletLogUUID.String(),
	).First(&walletLog).Error; err != nil {
		return nil, err
	} else {
		return &walletLog, nil
	}
}

func NewWalletRepository(db *gorm.DB) IWalletsRepository {
	return &WalletsRepository{
		db: db,
	}
}
