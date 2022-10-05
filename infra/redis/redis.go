package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/dougefr/product-scrapper/domain/contract/cache"
	"github.com/dougefr/product-scrapper/domain/contract/env"
	"github.com/dougefr/product-scrapper/domain/contract/logger"
	redis2 "github.com/go-redis/redis/v8"
)

type (
	redis[T any] struct {
		logger logger.Logger
		client redis2.Client
	}
)

// NewRedis cria um novo cache.Cache[t]
func NewRedis[T any](
	env env.Env,
	logger logger.Logger,
) cache.Cache[T] {
	client := redis2.NewClient(&redis2.Options{
		Addr: env.CacheAddr(),
	})

	return redis[T]{
		logger: logger,
		client: *client,
	}
}

func (r redis[T]) Set(
	ctx context.Context,
	key string,
	value T,
	expiration time.Duration,
) (err error) {

	ba, err := json.Marshal(value)
	if err != nil {
		r.logger.Error(ctx, "error when marshaling json", logger.Body{
			"err":   err,
			"key":   key,
			"value": value,
		})

		err = fmt.Errorf("error when marshaling json: %w", err)

		return
	}

	result := r.client.SetEX(ctx, key, string(ba), expiration)
	err = result.Err()
	if err != nil {
		r.logger.Error(ctx, "error when setting data in redis", logger.Body{
			"err":   err,
			"key":   key,
			"value": value,
		})

		err = fmt.Errorf("error when setting data in redis: %w", err)

		return
	}

	return
}

func (r redis[T]) Get(ctx context.Context, key string) (value *T, err error) {
	result := r.client.Get(ctx, key)
	err = result.Err()
	if err != nil {
		r.logger.Error(ctx, "error when getting data in Redis", logger.Body{
			"err": err,
			"key": key,
		})

		err = fmt.Errorf("error when getting data in Redis: %w", err)

		return
	}

	err = json.Unmarshal([]byte(result.Val()), &value)
	if err != nil {
		r.logger.Error(ctx, "error when unmarshalling json", logger.Body{
			"err":   err,
			"key":   key,
			"value": value,
		})

		err = fmt.Errorf("error when unmarshaling json: %w", err)

		return
	}

	return
}
