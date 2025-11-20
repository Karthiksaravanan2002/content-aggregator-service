package features

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// FeatureStrategy processes normalized provider content.
// Example features: trending, mylist, continue watching.
type FeatureStrategy interface {
	Name() string
	Execute(ctx context.Context, raw []models.ContentItem, req models.ProviderRequest) []models.ContentItem
}
