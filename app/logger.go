package app

import (
	"github.com/gin-contrib/logger"
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/thoas/go-funk"
)

type LoggerMiddleware gin.HandlerFunc

func NewLoggerMiddleware() LoggerMiddleware {
	zerolog.SetGlobalLevel(funk.ShortIf(
		gin.IsDebugging(),
		zerolog.DebugLevel,
		zerolog.InfoLevel,
	).(zerolog.Level))
	return LoggerMiddleware(logger.SetLogger())
}
