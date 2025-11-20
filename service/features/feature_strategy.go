package features

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// FeatureStrategy processes raw provider content.
// Example features: trending, mylist, continue watching.
type FeatureStrategy interface {
	Name() string
	Apply(ctx context.Context, raw map[string][]models.ContentItem) ([]models.ContentItem, error)
}
