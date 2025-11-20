package providers

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// FeatureStrategy processes raw provider content.
// Example features: trending, mylist, continue watching.
type FeatureStrategy interface {
	Apply(ctx context.Context, raw []models.ContentItem) ([]models.ContentItem, error)
}
