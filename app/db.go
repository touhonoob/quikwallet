package app

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"log"
	"os"
	"time"
	"github.com/go-redis/redis/v8"
)

func NewDb() *gorm.DB {
	if db, err := gorm.Open(mysql.Open(os.Getenv("QUIKWALLET_SQL_DSN")), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,   // Slow SQL threshold
				LogLevel:      logger.Silent, // Log level
				Colorful:      false,         // Disable color
			},
		),
	}); err != nil {
		panic(err)
	} else {
		return db
	}
}

func NewRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", os.Getenv("QUIKWALLET_REDIS_HOST"), os.Getenv("QUIKWALLET_REDIS_PORT")),
	})
}