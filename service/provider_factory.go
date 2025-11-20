package service

import (
	"strings"

	"dev.azure.com/daimler-mic/content-aggregator/service/cache"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers/twitch"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers/youtube"
	"go.uber.org/zap"
)

type ProviderFactory struct {
	props       props.ProvidersConfig
	cacheConfig props.CacheConfig
	redisCache  cache.Cache
	logger      *zap.Logger
}

func (f *ProviderFactory) GetProvider(name string) providers.ProviderStrategy {

	var provider providers.ProviderStrategy

	switch strings.ToLower(name) {
	case "youtube":
		provider = youtube.NewYouTubeStrategy(f.props.YouTube, f.logger)
	case "twitch":
		provider = twitch.NewTwitchStrategy(f.props.Twitch, f.logger)

	default:
		f.logger.Warn("unknown provider requested", zap.String("provider", name))
		return nil
	}

	cached := cache.NewCacheDecorator(
		provider,
		f.redisCache,
		f.logger,
		f.cacheConfig.TTLSeconds,
	)

	return cached
}

func NewProviderFactory(providerProps props.ProvidersConfig, cacheProps props.CacheConfig, logger *zap.Logger) *ProviderFactory {

	redis := cache.NewRedisCache(cacheProps)

	return &ProviderFactory{
		props:       providerProps,
		cacheConfig: cacheProps,
		redisCache:  redis,
		logger:      logger.With(zap.String("component", "provider-factory")),
	}
}
