package memory

import (
	"context"
	"time"

	"github.com/patrickmn/go-cache"
)

type Cacher struct {
	*cache.Cache
}

func New() Cacher {
	cacher := cache.New(5*time.Minute, 10*time.Minute)
	return Cacher{cacher}
}

func (c Cacher) Get(_ context.Context, key string) (any, bool) {
	return c.Cache.Get(key)
}

func (c Cacher) Set(_ context.Context, key string, x any, d time.Duration) {
	c.Cache.Set(key, x, d)
}

func (c Cacher) Clear(_ context.Context) error {
	c.Cache.Flush()
	return nil
}
