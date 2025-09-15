package redis

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Cacher struct {
	rdb *redis.Client
}

type Option struct {
	Addr     string
	Username string
	Password string
	DB       int
}

// New предоставляет инициализацию с возможной ошибкой (например, ping по redis)
func New(opt *Option) (*Cacher, error) {

	rdb := redis.NewClient(&redis.Options{
		Addr:     opt.Addr,
		Username: opt.Username,
		Password: opt.Password,
		DB:       opt.DB,
	})

	// Проверяем соединение
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &Cacher{rdb: rdb}, nil
}

func (c *Cacher) Get(ctx context.Context, key string, dest any) (bool, error) {

	val, err := c.rdb.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return false, nil
	}
	if err != nil {
		return false, err
	}

	if err = json.Unmarshal(val, dest); err != nil {
		return false, err
	}
	return true, nil
}

func (c *Cacher) Set(ctx context.Context, key string, value any, expire time.Duration) error {

	bytes, err := json.Marshal(value)
	if err != nil {
		return err
	}
	if err := c.rdb.Set(ctx, key, bytes, expire).Err(); err != nil {
		return err
	}
	return nil
}

func (c *Cacher) Clear(ctx context.Context) error {
	return c.rdb.FlushDB(ctx).Err()
}

func (c *Cacher) ClearByPrefix(ctx context.Context, prefix string) error {
	// TODO метод, который очищает кэш по ID проекта
	// Сканируем ключи по паттерну "prefix*"
	var cursor uint64
	pattern := prefix + "*"
	for {
		keys, nextCursor, err := c.rdb.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return err
		}
		if len(keys) > 0 {
			if err := c.rdb.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		cursor = nextCursor
		if cursor == 0 {
			break
		}
	}
	return nil
}

func (c *Cacher) Incr(ctx context.Context, key string, expire time.Duration) (int64, error) {

	// INCR делает инкремент или создает ключ, если его не было (в Redis по умолчанию 0=>1)
	val, err := c.rdb.Incr(ctx, key).Result()
	if err != nil {
		return 0, err
	}

	// Если только что создали ключ, надо явно задать expire (т.к. INCR для нового не ставит срок жизни)
	if val == 1 && expire > 0 {
		c.rdb.Expire(ctx, key, expire)
	}

	return val, nil
}
