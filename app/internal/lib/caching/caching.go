package caching

import (
	"context"
	"time"
)

type Cacher interface {
	Get(ctx context.Context, key string, dest any) (bool, error)
	Set(ctx context.Context, key string, value any, expire time.Duration) error
	Incr(ctx context.Context, key string, expire time.Duration) (int64, error)
	Clear(ctx context.Context) error
	ClearByPrefix(ctx context.Context, prefix string) error
}

type Service struct {
	Cacher
}

func New(c Cacher) Service {
	return Service{
		Cacher: c,
	}
}
