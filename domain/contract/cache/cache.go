package cache

import (
	"context"
	"time"
)

type (
	// Cache interface que representa uma infraestrutura de cache
	Cache[T any] interface {
		Set(ctx context.Context, key string, value T, expiration time.Duration) error
		Get(ctx context.Context, key string) (*T, error)
	}
)
