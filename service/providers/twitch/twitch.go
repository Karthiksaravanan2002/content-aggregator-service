package twitch

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

type TwitchStrategy struct {
	props  props.TwitchConfig
	logger *zap.Logger
}

func NewTwitchStrategy(props props.TwitchConfig, logger *zap.Logger) *TwitchStrategy {
	return &TwitchStrategy{
		props:  props,
		logger: logger.With(zap.String("provider", "twitch")),
	}
}

func (t *TwitchStrategy) FetchContent(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	return map[string][]models.ContentItem{}, nil
}
