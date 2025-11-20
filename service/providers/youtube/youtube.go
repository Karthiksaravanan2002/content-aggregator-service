package youtube

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
)

type YouTubeStrategy struct {
	props  props.YoutubeConfig
	logger *zap.Logger
}

func NewYouTubeStrategy(props props.YoutubeConfig, logger *zap.Logger) *YouTubeStrategy {
	return &YouTubeStrategy{
		props:  props,
		logger: logger.With(zap.String("provider", "youtube")),
	}
}

func (y *YouTubeStrategy) FetchContent(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	return map[string][]models.ContentItem{}, nil
}
