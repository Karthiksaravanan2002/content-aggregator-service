package cache

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	"go.uber.org/zap"
)

type CacheDecorator struct {
	inner  providers.ProviderStrategy
	cache  Cache
	logger *zap.Logger
	cfg props.CacheConfig
}

func NewCacheDecorator(inner providers.ProviderStrategy, cache Cache, logger *zap.Logger, cfg props.CacheConfig) providers.ProviderStrategy {
	return &CacheDecorator{
		inner:  inner,
		cache:  cache,
		logger: logger,
		cfg: cfg,
	}
}

// Feature-level caching
func (d *CacheDecorator) FetchFeatureRaw(ctx context.Context, req models.ProviderRequest, feature string) (map[string][]models.ContentItem, error) {


// search cache

	raw, err := d.inner.FetchFeatureRaw(ctx, req, feature)
	if err != nil {
		return nil, err
	}

//set cache

	return raw, nil
}
