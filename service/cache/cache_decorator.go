package cache

import (
	"context"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	"go.uber.org/zap"
)

type CacheDecorator struct {
	Inner  providers.ProviderStrategy
	Cache  Cache
	Logger *zap.Logger
	TTL    time.Duration
}

func (d *CacheDecorator) FetchContent(
	ctx context.Context,
	req models.ProviderRequest,
) (map[string][]models.ContentItem, error) {

	key := d.buildKey(req)

	// Try to hit cache
	if cached, ok := d.Cache.Get(key); ok {
		d.Logger.Info("cache hit", zap.String("key", key))
		return d.deserialize(cached), nil
	}

	d.Logger.Info("cache miss", zap.String("key", key))

	// Call actual provider
	data, err := d.Inner.FetchContent(ctx, req)
	if err != nil {
		return nil, err
	}

	// Save to cache
	d.Cache.Set(key, d.serialize(data), d.TTL)

	return data, nil
}

func (d *CacheDecorator) buildKey(req models.ProviderRequest) string {
	key := req.Provider
	for _, f := range req.Functionality {
		key += ":" + f
	}
	return key
}

func (d *CacheDecorator) serialize(
	in map[string][]models.ContentItem,
) map[string][]interface{} {

	return map[string][]interface{}{}
}

func (d *CacheDecorator) deserialize(
	in map[string][]interface{},
) map[string][]models.ContentItem {
	return map[string][]models.ContentItem{}
}

func NewCacheDecorator(
	inner providers.ProviderStrategy,
	cache Cache,
	logger *zap.Logger,
	ttlSeconds int,
) *CacheDecorator {

	return &CacheDecorator{
		Inner:  inner,
		Cache:  cache,
		Logger: logger.With(zap.String("layer", "cache")),
		TTL:    time.Duration(ttlSeconds) * time.Second,
	}
}
