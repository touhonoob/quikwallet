package app

import "github.com/gin-gonic/gin"

type IQuikWalletApp interface {
	Run() error
	Router() gin.IRouter
}