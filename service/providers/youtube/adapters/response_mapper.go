package adapters

import (
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	helper "dev.azure.com/daimler-mic/content-aggregator/service/providers/youtube/helpers"
	"google.golang.org/api/youtube/v3"
)

// MapYouTubeResponse converts raw YouTube API response from FetchContent
// into map[string][]models.ContentItem as required by ProviderStrategy.
func MapYouTubeResponse(v *youtube.Video) *models.ContentItem {
	sn := v.Snippet
	st := v.Statistics
	cd := v.ContentDetails

	// Parse RFC3339 published time safely
	var publishedAtPtr *time.Time
	if sn != nil && sn.PublishedAt != "" {
		if t, err := time.Parse(time.RFC3339, sn.PublishedAt); err == nil {
			publishedAtPtr = &t
		}
	}

	return &models.ContentItem{
		ID:           v.Id,
		Provider:     "youtube",
		ContentType:  models.ContentTypeVideo,
		Title:        helper.DefaultString(sn.Title, ""),
		Description:  helper.DefaultString(sn.Description, ""),
		PublishedAt:  models.PublishedAt{
			Timestamp: *publishedAtPtr,
		},
		ThumbnailURL: helper.ExtractBestThumbnailFromYT(sn.Thumbnails),
		ContentURL:   helper.MakeVideoURL(v.Id),

		ViewCount:       models.ViewCount{
			Value: helper.ToInt64(st.ViewCount),
		},

		ChannelID:  helper.DefaultString(sn.ChannelId, ""),
		Channel:    helper.DefaultString(sn.ChannelTitle, ""),
		ChannelURL: helper.MakeChannelURL(sn.ChannelId),

		Extras: &models.Extras{
			DurationSeconds: &models.Duration{
				Seconds: helper.ParseISODuration(helper.DefaultString(cd.Duration, "")),
			},
		},
	}
}

// Helper function to safely get string values
func getString(m map[string]any, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
