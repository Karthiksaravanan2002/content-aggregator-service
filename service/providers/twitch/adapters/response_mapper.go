package adapters

import (
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	helper "dev.azure.com/daimler-mic/content-aggregator/service/providers/twitch/helpers"
)

func MapTwitchStream(m map[string]any) *models.ContentItem {
	startedAt := helper.MustParseTime(helper.GetString(m, "started_at"))

	return &models.ContentItem{
		ID:          helper.GetString(m, "id"),
		Provider:    "twitch",
		ContentType: models.ContentTypeLive,

		Title:       helper.GetString(m, "title"),
		Description: "",

		PublishedAt:  models.PublishedAt{Timestamp: startedAt},
		ThumbnailURL: helper.ReplaceThumbnail(helper.GetString(m, "thumbnail_url")),
		ContentURL:   helper.MakeStreamURL(helper.GetString(m, "user_login")),

		ViewCount: models.ViewCount{
			Value: helper.GetInt64(m, "viewer_count"),
		},

		ChannelID:  helper.GetString(m, "user_id"),
		Channel:    helper.GetString(m, "user_name"),
		ChannelURL: helper.MakeChannelURL(helper.GetString(m, "user_id")),

		Extras: nil,
	}
}
