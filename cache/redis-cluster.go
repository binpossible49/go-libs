package cache

import (
	"context"
	"encoding/json"
	"time"

	"github.com/binpossible49/go-libs/opentracing/jaeger"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go/ext"
)

func initRedisCluster(addrs []string) (*redis.ClusterClient, error) {
	clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs: addrs,
	})
	_, err := clusterClient.Ping().Result()
	return clusterClient, err
}

type clusterRedisHelper struct {
	clusterClient *redis.ClusterClient
}

func (h *clusterRedisHelper) Exists(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Exists", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	indicator, err := h.clusterClient.Exists(key).Result()
	if err != nil {
		return err
	}
	if indicator == 0 {
		return redis.Nil
	}
	return nil
}

func (h *clusterRedisHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Get", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := h.clusterClient.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Set", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	_, err = h.clusterClient.Set(key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Del(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Del", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.clusterClient.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *clusterRedisHelper) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.clusterRedisHelper/Expire", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.clusterClient.Expire(key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}
