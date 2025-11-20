package youtube

import (
	"context"
	"fmt"
	"net/http"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

type YouTubeStrategy struct {
	props  props.YoutubeConfig
	logger *zap.Logger
}

func NewYouTubeStrategy(props props.YoutubeConfig, logger *zap.Logger) *YouTubeStrategy {
	return &YouTubeStrategy{
		props:  props,
		logger: logger.With(zap.String("provider", "youtube")),
	}
}

func (y *YouTubeStrategy) FetchContent(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	client := &http.Client{
		Transport: &transport.APIKey{Key: y.props.ApiKey},
	}

	// Create YouTube service instance
	svc, err := youtube.New(client)
	if err != nil {
		y.logger.Error("failed to create YouTube service", zap.Error(err))
		return nil, fmt.Errorf("youtube service init error: %w", err)
	}

	// Query parameters from ProviderRequest (example: search query)
	query := y.props.SearchQuery
	if query == "" {
		query = "trending"
	}

	// Call YouTube Search API
	searchCall := svc.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(10)

	resp, err := searchCall.Context(ctx).Do()
	if err != nil {
		y.logger.Error("youtube search failed", zap.Error(err))
		return nil, fmt.Errorf("youtube api error: %w", err)
	}

	// Prepare clean output
	output := make([]map[string]any, 0)

	for _, item := range resp.Items {
		if item.Id.Kind != "youtube#video" {
			continue
		}

		videoData := map[string]any{
			"videoId":     item.Id.VideoId,
			"title":       item.Snippet.Title,
			"description": item.Snippet.Description,
			"channel":     item.Snippet.ChannelTitle,
			"publishedAt": item.Snippet.PublishedAt,
			"thumbnails":  item.Snippet.Thumbnails.Default.Url,
		}

		output = append(output, videoData)
	}

	rawResp := map[string]any{
		"provider": "youtube",
		"results":  output,
	}

	return MapYouTubeResponse(rawResp)

}
