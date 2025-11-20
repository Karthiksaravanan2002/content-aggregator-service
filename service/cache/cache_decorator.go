package cache

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
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
func (d *CacheDecorator) FetchFeature(ctx context.Context,req *models.ProviderRequest,feature string) ([]*models.ContentItem, errors.AppError) {

// search cache
	items, err := d.inner.FetchFeature(ctx, req, feature)
	if err != nil {
		// Do not cache errors
		return nil, err
	}

	// 2 — Store in cache

	return items, nil
}
