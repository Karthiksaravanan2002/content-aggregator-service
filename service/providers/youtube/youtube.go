package youtube

import (
	"context"
	"errors"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"go.uber.org/zap"
	"google.golang.org/api/option"
	youtube "google.golang.org/api/youtube/v3"
)

const query = "trending"

type YouTubeProvider struct {
	cfg    props.YouTubeConfig
	logger *zap.Logger
	responseStore map[string]interface{}
}

func NewYouTubeProvider(cfg props.YouTubeConfig, logger *zap.Logger) *YouTubeProvider {
	return &YouTubeProvider{
		cfg:    cfg,
		logger: logger.With(zap.String("provider", "youtube")),
		responseStore: make(map[string]interface{}),
	}
}

func (yt *YouTubeProvider) FetchFeatureRaw(ctx context.Context, req models.ProviderRequest, feature string,) (map[string][]models.ContentItem, error) {
	switch feature {
	case "trending":
		return yt.callTrendingAPI(ctx, req)

	case "continue_watching":
		return yt.callContinueWatchingAPI(ctx, req)

	case "mylist":
		return yt.callMyListAPI(ctx, req)

	default:
		return nil, errors.New("feature not supported by youtube")
	}
}

// Skeleton implementations:
func (yt *YouTubeProvider) callTrendingAPI(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {


	// Create YouTube service instance
	svc, err := youtube.NewService(ctx, option.WithAPIKey(yt.cfg.ApiKey))
	if err != nil {
		yt.logger.Error("failed to create YouTube service", zap.Error(err))
		return nil, err
	}


	// Call YouTube Search API
	searchCall := svc.Search.List([]string{"id", "snippet"}).
		Q(query).
		MaxResults(10)

	resp, err := searchCall.Context(ctx).Do()
	if err != nil {
		yt.logger.Error("youtube search failed", zap.Error(err))
		return nil, err
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

func (yt *YouTubeProvider) callContinueWatchingAPI(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	return map[string][]models.ContentItem{}, nil
}

func (yt *YouTubeProvider) callMyListAPI(ctx context.Context, req models.ProviderRequest) (map[string][]models.ContentItem, error) {
	return map[string][]models.ContentItem{}, nil
}
