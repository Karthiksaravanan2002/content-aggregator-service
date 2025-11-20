package providers

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// ProviderStrategy defines a provider (YouTube, Twitch, etc.)
type ProviderStrategy interface {
	FetchContent(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error)
}
