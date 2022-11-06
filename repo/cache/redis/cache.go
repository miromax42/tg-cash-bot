package redis

import (
	"context"
	"time"

	"github.com/cockroachdb/errors"
	rcache "github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"gitlab.ozon.dev/miromaxxs/telegram-bot/repo/cache"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util"
	"gitlab.ozon.dev/miromaxxs/telegram-bot/util/metrics"
)

type Cache struct {
	client *redis.Client
	cch    *rcache.Cache

	objectTTL time.Duration
}

func NewCache(ctx context.Context, cfg util.ConfigCache) (*Cache, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.SocketAddr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	err := rdb.Ping(ctx).Err()
	if err != nil {
		return nil, errors.Wrapf(err, "ping redis on %q", cfg.Redis.SocketAddr)
	}

	rcch := rcache.New(&rcache.Options{
		Redis:      rdb,
		LocalCache: rcache.NewTinyLFU(cfg.ObjectsCount, cfg.TTL),
	})

	return &Cache{
		client:    rdb,
		cch:       rcch,
		objectTTL: cfg.ObjectTTL,
	}, nil
}

func (c *Cache) Get(ctx context.Context, key cache.Token, value interface{}) error {
	err := c.cch.Get(ctx, string(key), value)
	if err != nil {
		if errors.Is(err, rcache.ErrCacheMiss) {
			metrics.CacheMissCounter.Add(1)

			return errors.Wrapf(cache.ErrMiss, "get %q", string(key))
		}

		return errors.Wrapf(err, "get %q", string(key))
	}

	metrics.CacheHitCounter.Add(1)

	return nil
}

func (c *Cache) Set(ctx context.Context, key cache.Token, value interface{}) error {
	err := c.cch.Set(&rcache.Item{
		Ctx:   ctx,
		Key:   string(key),
		Value: value,
		TTL:   c.objectTTL,
	})

	return errors.Wrapf(err, "set %q", string(key))
}

func (c *Cache) Del(ctx context.Context, key cache.Token) error {
	err := c.cch.Delete(ctx, string(key))

	return errors.Wrapf(err, "delete %q", string(key))
}
