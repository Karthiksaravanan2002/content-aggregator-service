package youtube

import (
	"context"

	stdErr "errors"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	yxadapters "dev.azure.com/daimler-mic/content-aggregator/service/providers/youtube/adapters"
	yxfeatures "dev.azure.com/daimler-mic/content-aggregator/service/providers/youtube/features"
	"go.uber.org/zap"
)

const query = "trending"

type YouTubeProvider struct {
	cfg    props.YouTubeConfig
	logger *zap.Logger
	responseStore map[string]interface{}
	adapters map[string]FeatureAdapter
	features map[string]providers.FeatureStrategy    
}

type FeatureAdapter func(ctx context.Context, req *models.ProviderRequest) ([]*models.ContentItem, errors.AppError)


func (yt *YouTubeProvider) FetchFeature(ctx context.Context, req *models.ProviderRequest, feature string,) ([]*models.ContentItem, errors.AppError) {
    adapter, ok := yt.adapters[feature]
    if !ok {
        return nil, errors.BadRequest(stdErr.New("Feature Unsupported"))
    }

    raw, err := adapter(ctx, req)
    if err != nil {
        return nil, err
    }

    feat, ok := yt.features[feature]
    if !ok {
        return nil, err
    }

    return feat.Apply(ctx, raw)
}

func (yt *YouTubeProvider) GetFeature(feature string) providers.FeatureStrategy {
    return yt.features[feature]
}

func NewYouTubeProvider(cfg props.YouTubeConfig, logger *zap.Logger) *YouTubeProvider {
	p:= &YouTubeProvider{
		cfg:    cfg,
		logger: logger.With(zap.String("provider", "youtube")),
		responseStore: make(map[string]interface{}),
		adapters: make(map[string]FeatureAdapter),
		features: make(map[string]providers.FeatureStrategy),
	}
     p.registerAdapters()
	   p.registerFeatures()

		 return p
}


func (yt *YouTubeProvider) registerAdapters() {
	yt.adapters["trending"] = yxadapters.CallTrendingAPI(yt.cfg, yt.logger)
	// yt.adapters["continue_watching"] = yxadapters.CallContinueWatchingAPI(yt.cfg, yt.logger)
	// yt.adapters["mylist"] = yxadapters.CallMyListAPI(yt.cfg, yt.logger)
}

func (yt *YouTubeProvider) registerFeatures() {
	yt.features["trending"] = yxfeatures.NewYouTubeTrendingFeature(yt.cfg, yt.logger)
	// yt.features["continue_watching"] = yxfeatures.NewYouTubeContinueFeature()
	// yt.features["mylist"] = yxfeatures.NewYouTubeMyListFeature()
}
