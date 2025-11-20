package providers

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// ProviderStrategy defines a provider (YouTube, Twitch, etc.)
type ProviderStrategy interface {
	FetchFeatureRaw(ctx context.Context, req models.ProviderRequest, feature string) (map[string][]models.ContentItem, error)
}

