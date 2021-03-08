package apiv1wallets

import (
	"context"
	"fmt"
	redisCache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

type WalletsCache struct{
	redisClient *redis.Client
	redisCache *redisCache.Cache
	repository IWalletsRepository
}

type GetWalletBalanceCache struct {
	Balance string
}

func (walletCache *WalletsCache) GetWalletBalance(walletUUID uuid.UUID) (decimal.Decimal, error) {
	var getWalletBalanceCache GetWalletBalanceCache
	if err := walletCache.redisCache.Once(&redisCache.Item{
		Key: walletCache.getWalletBalanceCacheKey(walletUUID),
		Value: &getWalletBalanceCache,
		Do: func(item *redisCache.Item) (interface{}, error) {
			if balance, err := walletCache.repository.GetBalance(walletUUID); err != nil {
				return nil, err
			} else {
				return &GetWalletBalanceCache{Balance: balance.String()}, nil
			}
		},
	}); err != nil {
		return decimal.Zero, err
	} else {
		return decimal.NewFromString(getWalletBalanceCache.Balance)
	}
}

func (walletCache *WalletsCache) InvalidateWalletBalance(walletUUID uuid.UUID) error {
	return walletCache.redisCache.Delete(context.Background(), walletCache.getWalletBalanceCacheKey(walletUUID))
}

func (WalletsCache) getWalletBalanceCacheKey(walletUUID uuid.UUID) string {
	return fmt.Sprintf("GetWalletBalance:%s", walletUUID.String())
}

func NewWalletsCache(redisClient *redis.Client, repository IWalletsRepository) IWalletsCache {
	return &WalletsCache{
		redisClient: redisClient,
		redisCache: redisCache.New(&redisCache.Options{
			Redis:      redisClient,
		}),
		repository: repository,
	}
}