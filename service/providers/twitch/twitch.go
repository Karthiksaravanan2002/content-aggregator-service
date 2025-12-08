package twitch

import (
	"context"
	stdErr "errors"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	txadapters "dev.azure.com/daimler-mic/content-aggregator/service/providers/twitch/adapters"
	txfeatures "dev.azure.com/daimler-mic/content-aggregator/service/providers/twitch/features"

	"go.uber.org/zap"
)

type FeatureAdapter func(ctx context.Context, req *models.ProviderRequest) ([]*models.ContentItem, errors.AppError)

type TwitchProvider struct {
	cfg           props.TwitchConfig
	logger        *zap.Logger
	responseStore map[string]interface{}
	adapters      map[string]FeatureAdapter
	features      map[string]providers.FeatureStrategy
}

func NewTwitchProvider(cfg props.TwitchConfig, logger *zap.Logger) *TwitchProvider {
	p := &TwitchProvider{
		cfg:           cfg,
		logger:        logger.With(zap.String("provider", "twitch")),
		responseStore: make(map[string]interface{}),
		adapters:      make(map[string]FeatureAdapter),
		features:      make(map[string]providers.FeatureStrategy),
	}

	p.registerAdapters()
	p.registerFeatures()
	return p
}

func (tp *TwitchProvider) FetchFeature(ctx context.Context, req *models.ProviderRequest, feature string) ([]*models.ContentItem, errors.AppError) {
	adapter, ok := tp.adapters[feature]
	if !ok {
		return nil, errors.BadRequest(stdErr.New("Feature Unsupported"))
	}

	raw, err := adapter(ctx, req)
	if err != nil {
		return nil, err
	}

	feat, ok := tp.features[feature]
	if !ok {
		return nil, err
	}

	return feat.Apply(ctx, raw)
}

func (tp *TwitchProvider) GetFeature(feature string) providers.FeatureStrategy {
	return tp.features[feature]
}

func (tp *TwitchProvider) registerAdapters() {
	tp.adapters["trending"] = txadapters.CallTrendingStreams(tp.cfg, tp.logger)
}

func (tp *TwitchProvider) registerFeatures() {
	tp.features["trending"] = txfeatures.NewTwitchTrendingFeature(tp.cfg, tp.logger)
}
