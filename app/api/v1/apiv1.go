package apiv1

import (
	"github.com/gin-gonic/gin"
	"github.com/touhonoob/quikwallet/app"
	"time"
)

type ApiV1 struct {
	routerGroup *gin.RouterGroup
}

type APIError struct {
	ErrorCode    int
	ErrorMessage string
	CreatedAt    time.Time
}

func (api *ApiV1) Router() gin.IRouter {
	return api.routerGroup
}

func NewApiV1(app app.IQuikWalletApp) IApiV1 {
	return &ApiV1{
		routerGroup: app.Router().Group("/api/v1"),
	}
}
