package twitch

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

type TwitchProvider struct {
	cfg    props.TwitchConfig
	logger *zap.Logger
	responseStore map[string]interface{}
}

func NewTwitchProvider(cfg props.TwitchConfig, logger *zap.Logger) *TwitchProvider {
	return &TwitchProvider{
		cfg:    cfg,
		logger: logger.With(zap.String("provider", "twitch")),
		responseStore: make(map[string]interface{}),
	}
}

func (t *TwitchProvider) FetchFeatureRaw(ctx context.Context, req models.ProviderRequest, feature string) (map[string][]models.ContentItem, error) {
	// Twitch has ONE endpoint → reuse for all features
	raw, _ := t.fetchStreamsAPI(ctx, req)
	return raw, nil
}

func (t *TwitchProvider) fetchStreamsAPI(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	return map[string][]models.ContentItem{}, nil
}
