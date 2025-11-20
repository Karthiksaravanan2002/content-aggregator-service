package trending

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// TrendingFeature simply passes trending items (already normalized).
// Filtering or ranking logic can go here.
type TrendingFeature struct{}

func NewTrendingFeature() *TrendingFeature {
	return &TrendingFeature{}
}

func (f *TrendingFeature) Name() string {
	return "trending"
}

func (f *TrendingFeature) Execute(
	ctx context.Context,
	items []models.ContentItem,
	req models.ProviderRequest,
) []models.ContentItem {

	return items
}
