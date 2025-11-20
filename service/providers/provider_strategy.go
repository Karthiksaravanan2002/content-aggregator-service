package providers

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// ProviderStrategy defines a provider (YouTube, Twitch, etc.)
type ProviderStrategy interface {
	FetchFeature(ctx context.Context, req *models.ProviderRequest, feature string) ([]*models.ContentItem, errors.AppError)
}

