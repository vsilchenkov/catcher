package memory

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/cockroachdb/errors"
	"github.com/patrickmn/go-cache"
)

var ErrWrongType = errors.New("value is not of the expected type")

type Cacher struct {
	*cache.Cache
}

func New() *Cacher {
	cacher := cache.New(5*time.Minute, 10*time.Minute)
	return &Cacher{cacher}
}

func (c *Cacher) Get(_ context.Context, key string, dest any) (bool, error) {

	val, found := c.Cache.Get(key)
	if !found {
		return false, nil
	}

	data, ok := val.([]byte)
	if !ok {
		return false, ErrWrongType
	}
	if err := json.Unmarshal(data, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Cacher) Set(_ context.Context, key string, x any, expire time.Duration) error {
	bytes, err := json.Marshal(x)
	if err != nil {
		return err
	}
	c.Cache.Set(key, bytes, expire)
	return nil
}

func (c *Cacher) Clear(_ context.Context) error {
	c.Cache.Flush()
	return nil
}

func (c *Cacher) ClearByPrefix(_ context.Context, prefix string) error {
    for k := range c.Cache.Items() {
        if strings.HasPrefix(k, prefix) {
            c.Cache.Delete(k)
        }
    }
    return nil
}

func (c *Cacher) Incr(ctx context.Context, key string, expire time.Duration) (int64, error) {

	var val int64
	c.Get(ctx, key, &val);

	val++
	c.Set(ctx, key, val, expire)

	return val, nil
}
