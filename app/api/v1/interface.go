package apiv1

import "github.com/gin-gonic/gin"

type IApiV1 interface {
	Router() gin.IRouter
}
