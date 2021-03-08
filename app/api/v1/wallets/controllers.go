package apiv1wallets

import (
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type WalletsController struct {
	repository IWalletsRepository
	cache IWalletsCache
	queue IQueue
}

func NewWalletsController(repository IWalletsRepository, cache IWalletsCache, queue IQueue) IWalletsController {
	return &WalletsController{
		repository: repository,
		cache: cache,
		queue: queue,
	}
}

func (controller *WalletsController) PostCredit(c *gin.Context) {
	var request PostCreditRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithError(400, err)
	} else if uuid, err := controller.getWalletUuidFromContext(c); err != nil {
		c.AbortWithError(400, errors.New("invalid wallet ID"))
	} else if credit, err := decimal.NewFromString(request.Credit); err != nil {
		c.AbortWithError(400, err)
	} else if creditLog, err := controller.repository.CreateCreditLog(
		uuid, credit,
	); err != nil {
		log.Error().Err(err).Msg("failed to create credit log")
		c.AbortWithError(500, errors.New("failed to create credit log"))
	} else if err := controller.queue.PublishNewWalletLogJob(uuid); err != nil {
		log.Error().Err(err).Msg("failed to publish new wallet log job")
		c.AbortWithError(500, errors.New("failed to process the request"))
	} else {
		c.JSON(202, &PostCreditResponse{
			WalletLogUuid: creditLog.Uuid,
		})
	}
}

type PostCreditRequest struct {
	Credit string `json:"credit" binding:"required,min=0,max=1000000"`
}

type PostCreditResponse struct{
	WalletLogUuid string `json:"wallet_log_uuid"`
}

func (controller *WalletsController) PostDebit(c *gin.Context) {
	var request PostDebitRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		c.AbortWithError(400, err)
	} else if uuid, err := controller.getWalletUuidFromContext(c); err != nil {
		c.AbortWithError(400, errors.New("invalid wallet ID"))
	} else if debit, err := decimal.NewFromString(request.Debit); err != nil {
		c.AbortWithError(400, err)
	} else if debitLog, err := controller.repository.CreateDebitLog(uuid, debit); err != nil {
		log.Error().Err(err).Msg("failed to create debit log")
		c.AbortWithError(500, errors.New("failed to create debit log"))
	} else if err := controller.queue.PublishNewWalletLogJob(uuid); err != nil {
		log.Error().Err(err).Msg("failed to publish new wallet log job")
		c.AbortWithError(500, errors.New("failed to process the request"))
	} else {
		c.JSON(202, &PostDebitResponse{
			WalletLogUuid: debitLog.Uuid,
		})
	}
}

type PostDebitRequest struct {
	Debit string `json:"debit" binding:"required,min=0,max=1000000"`
}

type PostDebitResponse struct{
	WalletLogUuid string `json:"wallet_log_uuid"`
}

func (controller *WalletsController) GetBalance(c *gin.Context) {
	if uuid, err := controller.getWalletUuidFromContext(c); err != nil {
		c.AbortWithError(400, errors.New("invalid wallet ID"))
	} else if _, err := controller.repository.GetWallet(uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithError(400, errors.New("invalid wallet ID"))
		} else {
			log.Error().Err(err).Msg("failed to get wallet")
			c.AbortWithError(500, errors.New("failed to get wallet"))
		}
	} else if balance, err := controller.cache.GetWalletBalance(uuid); err != nil {
		log.Error().Err(err).Msg("failed to get wallet balance")
		c.AbortWithError(500, errors.New("failed to get wallet balance"))
	} else {
		c.JSON(
			200, &GetBalanceResponse{
				Balance:  balance.Shift(-2).String(),
			},
		)
	}
}

type GetBalanceResponse struct {
	Balance  string `json:"balance"`
}

func (controller *WalletsController) GetWalletLog(c *gin.Context) {
	if uuid, err := controller.getWalletLogUuidFromContext(c); err != nil {
		c.AbortWithError(400, errors.New("invalid wallet log ID"))
	} else if walletLog, err := controller.repository.GetWalletLog(uuid); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.AbortWithError(400, errors.New("invalid wallet log ID"))
		} else {
			log.Error().Err(err).Msg("failed to get wallet log")
			c.AbortWithError(500, errors.New("failed to get wallet log"))
		}
	} else {
		c.JSON(200, walletLog)
	}
}

func (controller *WalletsController) getWalletUuidFromContext(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.Param("wallet_uuid"))
}

func (controller *WalletsController) getWalletLogUuidFromContext(c *gin.Context) (uuid.UUID, error) {
	return uuid.Parse(c.Param("wallet_log_uuid"))
}