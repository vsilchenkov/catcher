package caching

import (
	"context"
	"time"
)

type Cacher interface {
	Get(ctx context.Context, key string) (any, bool)
	Set(ctx context.Context, key string, x any, d time.Duration)
	Clear(ctx context.Context) error
}

type Service struct {
	Cacher
}

func New(c Cacher) Service {
	return Service{
		Cacher: c,
	}
}
