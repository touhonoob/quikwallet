package app

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"os"
)

type QuikWalletApp struct {
	ginEngine *gin.Engine
}

func (app *QuikWalletApp) Run() error {
	return app.ginEngine.Run(
		fmt.Sprintf(
			"%s:%s",
			os.Getenv("QUIKWALLET_HOST"),
			os.Getenv("QUIKWALLET_PORT"),
		),
	)
}

func (app *QuikWalletApp) Router() gin.IRouter {
	return app.ginEngine
}

func NewApp(loggerMiddleware LoggerMiddleware) IQuikWalletApp {
	app := &QuikWalletApp{
		ginEngine: gin.New(),
	}

	app.ginEngine.Use(gin.HandlerFunc(loggerMiddleware))
	app.ginEngine.Use(gin.Recovery())

	return app
}