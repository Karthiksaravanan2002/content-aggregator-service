package adapters

import (
	"context"
	"net/http"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	"google.golang.org/api/youtube/v3"
)

func CallTrendingAPI(cfg props.YouTubeConfig, logger *zap.Logger) func(context.Context, *models.ProviderRequest) ([]*models.ContentItem, errors.AppError) {
	return func(ctx context.Context, req *models.ProviderRequest) ([]*models.ContentItem, errors.AppError) {
		svc, err := youtube.NewService(ctx, option.WithAPIKey(cfg.ApiKey))
		if err != nil {
			logger.Error("failed to create YouTube service", zap.Error(err))
			return nil, errors.ProviderError(http.StatusBadGateway,err)
		}

		call := svc.Videos.List([]string{"id", "snippet", "statistics", "contentDetails"}).
			Chart("mostPopular").
			RegionCode("DE").
			MaxResults(20)

		resp, err := call.Context(ctx).Do()
		if err != nil {
			logger.Error("youtube trending API failed", zap.Error(err))
			return nil, errors.ProviderError(http.StatusBadGateway,err)
		}

		// Directly create final items WITHOUT intermediate map
		items := make([]*models.ContentItem, 0, len(resp.Items))

		for _, v := range resp.Items {
			if v.Snippet == nil {
				continue
			}
			items = append(items, MapYouTubeResponse(v))
		}

		return items,nil

	}
}
