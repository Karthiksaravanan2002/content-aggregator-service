package features

import (
	"context"
	"sort"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	"go.uber.org/zap"
)

type TwitchTrendingFeature struct {
	cfg    props.TwitchConfig
	logger *zap.Logger
}

// Constructor (same style as YouTube)
func NewTwitchTrendingFeature(cfg props.TwitchConfig, logger *zap.Logger) providers.FeatureStrategy {
	return &TwitchTrendingFeature{
		cfg:    cfg,
		logger: logger.With(zap.String("feature", "twitch_trending")),
	}
}

// Apply transforms raw stream data into final trending list
func (f *TwitchTrendingFeature) Apply(ctx context.Context, items []*models.ContentItem) ([]*models.ContentItem, errors.AppError) {

	if items == nil {
		return []*models.ContentItem{}, nil
	}

	// Sort by view count descending
	sort.Slice(items, func(i, j int) bool {
		return items[i].ViewCount.Value > items[j].ViewCount.Value
	})

	// Keep top N (20)
	if len(items) > 20 {
		items = items[:20]
	}

	// Normalization: ensure provider and content type set
	for _, item := range items {
		if item == nil {
			continue
		}
		item.Provider = "twitch"
		if item.ContentType == "" {
			item.ContentType = models.ContentTypeLive
		}
	}

	return items, nil
}
