package apiv1wallets

import (
	"github.com/gin-gonic/gin"
	"github.com/touhonoob/quikwallet/app/api/v1"
)

type ApiV1Wallets struct {
	routerGroup *gin.RouterGroup
}

func (api *ApiV1Wallets) Router() gin.IRouter {
	return api.routerGroup
}

func NewApiV1Wallets(apiV1 apiv1.IApiV1, controller IWalletsController) IApiV1Wallets {
	routerGroup := apiV1.Router().Group("/wallets")
	routerGroup.GET("/:wallet_uuid/balance", controller.GetBalance)
	routerGroup.POST("/:wallet_uuid/credit", controller.PostCredit)
	routerGroup.POST("/:wallet_uuid/debit", controller.PostDebit)
	routerGroup.GET("/:wallet_uuid/logs/:wallet_log_uuid", controller.GetWalletLog)

	return &ApiV1Wallets{
		routerGroup: routerGroup,
	}
}