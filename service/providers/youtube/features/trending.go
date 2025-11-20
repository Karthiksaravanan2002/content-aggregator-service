package features

import (
	"context"
	"sort"

	"go.uber.org/zap"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
)

type YouTubeTrendingFeature struct {
	logger *zap.Logger
	cfg    props.YouTubeConfig
}

func NewYouTubeTrendingFeature(cfg props.YouTubeConfig, logger *zap.Logger) providers.FeatureStrategy {
	return &YouTubeTrendingFeature{
		cfg:    cfg,
		logger: logger,
	}
}

func (f *YouTubeTrendingFeature) Apply(ctx context.Context, raw []models.ContentItem) ([]models.ContentItem, error) {

	sort.SliceStable(raw, func(i, j int) bool {
		return raw[i].ViewCount > raw[j].ViewCount
	})

	return raw, nil
}
