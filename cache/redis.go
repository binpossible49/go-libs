package cache

import (
	"context"
	"encoding/json"
	"reflect"
	"time"

	"github.com/binpossible49/go-libs/opentracing/jaeger"
	"github.com/go-redis/redis"
	"github.com/opentracing/opentracing-go/ext"
)

type redisHelper struct {
	client *redis.Client
}

func initRedis(addr string) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
		DB:   0,
	})
	_, err := client.Ping().Result()
	return client, err
}

func (h *redisHelper) Exists(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.redisHelper/Exists", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	indicator, err := h.client.Exists(key).Result()
	if err != nil {
		return err
	}
	if indicator == 0 {
		return redis.Nil
	}
	return nil
}

func (h *redisHelper) Get(ctx context.Context, key string, value interface{}) (err error) {
	span := jaeger.Start(ctx, ">helper.redisHelper/Get", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := h.client.Get(key).Result()
	if err != nil {
		return err
	}
	err = json.Unmarshal([]byte(data), &value)
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.redisHelper/Set", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	_, err = h.client.Set(key, string(data), expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Del(ctx context.Context, key string) (err error) {
	span := jaeger.Start(ctx, ">helper.redisHelper/Del", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.client.Del(key).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) Expire(ctx context.Context, key string, expiration time.Duration) (err error) {
	span := jaeger.Start(ctx, ">helper.redisHelper/Expire", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	_, err = h.client.Expire(key, expiration).Result()
	if err != nil {
		return err
	}
	return nil
}

func (h *redisHelper) GetInterface(ctx context.Context, key string, value interface{}) (interface{}, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/GetInterface", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()

	data, err := h.client.Get(key).Result()
	if err != nil {
		return nil, err
	}

	typeValue := reflect.TypeOf(value)
	kind := typeValue.Kind()

	var outData interface{}
	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		outData = reflect.New(typeValue).Interface()
	default:
		outData = reflect.Zero(typeValue).Interface()
	}
	err = json.Unmarshal([]byte(data), &outData)
	if err != nil {
		return nil, err
	}

	switch kind {
	case reflect.Ptr, reflect.Struct, reflect.Slice:
		return reflect.ValueOf(outData).Elem().Interface(), nil
	}
	var outValue interface{} = outData
	if reflect.TypeOf(outData).ConvertibleTo(typeValue) {
		outValueConverted := reflect.ValueOf(outData).Convert(typeValue)
		outValue = outValueConverted.Interface()
	}
	return outValue, nil
}

func (h *redisHelper) DelMulti(ctx context.Context, keys ...string) error {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/DelMulti", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	pipeline := h.client.TxPipeline()
	pipeline.Del(keys...)
	_, err = pipeline.Exec()
	return err
}

func (h *redisHelper) GetKeysByPattern(ctx context.Context, pattern string, cursor uint64, limit int64) ([]string, uint64, error) {
	var err error
	span := jaeger.Start(ctx, ">helper.redisHelper/GetKeysByPattern", ext.SpanKindRPCClient)
	defer func() {
		jaeger.Finish(span, err)
	}()
	return h.client.Scan(cursor, pattern, limit).Result()
}
