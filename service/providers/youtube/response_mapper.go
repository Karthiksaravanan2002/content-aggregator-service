package youtube

import (
	"fmt"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// MapYouTubeResponse converts raw YouTube API response from FetchContent
// into map[string][]models.ContentItem as required by ProviderStrategy.
func MapYouTubeResponse(raw map[string]any) (map[string][]models.ContentItem, error) {
	results, ok := raw["results"].([]map[string]any)
	if !ok {
		return nil, fmt.Errorf("invalid response format: missing 'results'")
	}

	contentItems := make([]models.ContentItem, 0)

	for _, item := range results {
		// Validate ID exists
		id, _ := item["videoId"].(string)
		if id == "" {
			continue
		}

		// Parse publishedAt timestamp
		var publishedTimePtr *time.Time
		if publishedStr, ok := item["publishedAt"].(string); ok {
			if t, err := time.Parse(time.RFC3339, publishedStr); err == nil {
				publishedTimePtr = &t
			}
		}

		// Handle thumbnails (from "thumbnails": string or nested struct)
		thumbnail := map[string]string{}
		if thumbURL, ok := item["thumbnails"].(string); ok {
			thumbnail["default"] = thumbURL
		}

		// Build ContentItem object
		contentItems = append(contentItems, models.ContentItem{
			ID:          id,
			Provider:    "youtube",
			Title:       getString(item, "title"),
			Description: getString(item, "description"),
			PublishedAt: publishedTimePtr,
			Channel:     getString(item, "channel"),
			Thumbnail:   thumbnail,
			Extra:       map[string]string{},
		})
	}

	return map[string][]models.ContentItem{
		"youtube": contentItems,
	}, nil
}

// Helper function to safely get string values
func getString(m map[string]any, key string) string {
	if val, ok := m[key].(string); ok {
		return val
	}
	return ""
}
