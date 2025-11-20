package trending

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

type TrendingFeature struct{}

func NewTrendingFeature() *TrendingFeature { return &TrendingFeature{} }
func (TrendingFeature) Name() string       { return "trending" }

func (TrendingFeature) Apply(ctx context.Context, raw map[string][]models.ContentItem) ([]models.ContentItem, error) {
return raw["youtube"], nil
}
