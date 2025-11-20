package continuewatching

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// ContinueWatchingFeature returns content user left unfinished.
type ContinueWatchingFeature struct{}

func NewContinueWatchingFeature() *ContinueWatchingFeature {
	return &ContinueWatchingFeature{}
}

func (f *ContinueWatchingFeature) Name() string {
	return "continue_watching"
}

func (f *ContinueWatchingFeature) Execute(
	ctx context.Context,
	items []models.ContentItem,
	req models.ProviderRequest,
) []models.ContentItem {

	return items
}
