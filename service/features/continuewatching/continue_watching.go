package continuewatching

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

type ContinueWatchingFeature struct{}

func NewContinueWatchingFeature() *ContinueWatchingFeature { return &ContinueWatchingFeature{} }
func (ContinueWatchingFeature) Name() string                { return "continue_watching" }

func (ContinueWatchingFeature) Apply(ctx context.Context, raw map[string][]models.ContentItem) ([]models.ContentItem, error) {
	return []models.ContentItem{}, nil
}
