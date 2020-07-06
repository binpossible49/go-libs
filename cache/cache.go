package cache

import (
	"context"
	"time"

	"go.uber.org/zap"
)

// CacheHelper is helper of Cache
type CacheHelper interface {
	Exists(ctx context.Context, key string) error
	Get(ctx context.Context, key string, value interface{}) error
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Del(ctx context.Context, key string) error
	Expire(ctx context.Context, key string, expiration time.Duration) error
}

// NewCacheHelper creates an instance
func NewCacheHelper(addrs []string) CacheHelper {
	if len(addrs) > 1 {
		clusterClient, err := initRedisCluster(addrs)
		if err != nil {
			zap.S().Panic("Failed to init redis cluster", zap.Error(err))
		}
		return &clusterRedisHelper{
			clusterClient: clusterClient,
		}
	}
	client, err := initRedis(addrs[0])
	if err != nil {
		zap.S().Panic("Failed to init redis", zap.Error(err))
	}
	return &redisHelper{
		client: client,
	}
}
