package features

import (
	"context"
	"sort"
	"time"

	"go.uber.org/zap"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	helper "dev.azure.com/daimler-mic/content-aggregator/service/providers/youtube/helpers"
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

func (f *YouTubeTrendingFeature) Apply(ctx context.Context, raw []*models.ContentItem) ([]*models.ContentItem, errors.AppError) {

  if len(raw) == 0 {
		f.logger.Info("no trending items returned from provider")
		return nil, nil
	}

	filtered := make([]*models.ContentItem, 0, len(raw))

	for _, it := range raw {

		// TV apps usually avoid showing Shorts in Trending
		if it.ContentType != models.ContentTypeVideo {
			continue
		}

		// Must have thumbnail
		if it.ThumbnailURL == "" {
			continue
		}

		// Must have title
		if it.Title == "" {
			continue
		}

		filtered = append(filtered, it)
	}
	if len(filtered) == 0 {
		f.logger.Warn("all trending items filtered out")
		return nil, nil
	}
	unique := make([]*models.ContentItem, 0, len(filtered))
	seen := make(map[string]bool)

	for _, it := range filtered {
		if !seen[it.ID] {
			seen[it.ID] = true
			unique = append(unique, it)
		}
	}

	// Sort by view count → Trending relevance
	sort.Slice(unique, func(i, j int) bool {
		return unique[i].ViewCount.Value > unique[j].ViewCount.Value
	})

	// Enrich relative timestamp and view formatting
	now := time.Now()
	for i := range unique {

		// Relative time: "2 days ago"
		if !unique[i].PublishedAt.Timestamp.IsZero() {
			unique[i].PublishedAt.Relative = helper.RelativeTime(now, unique[i].PublishedAt.Timestamp)
		}

		// Format view count
		unique[i].ViewCount.Display = helper.FormatViewCount(unique[i].ViewCount.Value)
	}

	f.logger.Info("youtube trending processed",
		zap.Int("inputCount", len(raw)),
		zap.Int("outputCount", len(unique)),
	)

	return unique, nil
}
