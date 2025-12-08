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


type ProviderFactory interface {
    GetProvider(name string) providers.ProviderStrategy
}

type providerFactory struct {
    props   props.ProvidersConfig
		cacheCfg props.CacheConfig
    cache   cache.Cache
    logger  *zap.Logger
    makers  map[string]ProviderMaker
}

type ProviderMaker func(logger *zap.Logger) providers.ProviderStrategy


func (f *providerFactory) Register(name string, maker ProviderMaker) {
    f.makers[strings.ToLower(name)] = maker
}


func (f *providerFactory) GetProvider(name string) providers.ProviderStrategy {
    maker, ok := f.makers[strings.ToLower(name)]
    if !ok {
        return nil
    }

    providerImpl := maker(f.logger)

    // Wrap with cache decorator
    return cache.NewCacheDecorator(
        providerImpl,
        f.cache,
        f.logger,
				f.cacheCfg,
    )
}


func NewProviderFactory(cfg props.ProvidersConfig, cacheCfg props.CacheConfig,cacheLayer cache.Cache, logger *zap.Logger,) ProviderFactory {

    f := &providerFactory{
        props:  cfg,
				cacheCfg: cacheCfg,
        cache:  cacheLayer,
        logger: logger,
        makers: make(map[string]ProviderMaker),
    }

    // Register providers here
    f.Register("youtube", func(logger *zap.Logger) providers.ProviderStrategy {
        return youtube.NewYouTubeProvider(cfg.YouTube, logger)
    })

    f.Register("twitch", func(logger *zap.Logger) providers.ProviderStrategy {
        return twitch.NewTwitchProvider(cfg.Twitch, logger)
    })

    return f
}

