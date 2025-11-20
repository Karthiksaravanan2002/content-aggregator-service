package models

import "time"

type ContentType string

const (
	// Standard long-form video content (YouTube videos, OTT extras)
	ContentTypeVideo ContentType = "video"

	// Short-form vertical content (YouTube Shorts, Reels-like content)
	ContentTypeShort ContentType = "short"

	// Very short extracted moment from a stream/video (Twitch Clips)
	ContentTypeClip ContentType = "clip"

	// A full livestream that is currently live
	ContentTypeLive ContentType = "live"

	// A scheduled livestream that has not started yet
	ContentTypeUpcoming ContentType = "upcoming"

	// Video on demand (Twitch past broadcasts, recorded streams)
	ContentTypeVOD ContentType = "vod"

	// Longer curated segment from a livestream (Twitch Highlights)
	ContentTypeHighlight ContentType = "highlight"

	// A full-length movie from OTT platforms
	ContentTypeMovie ContentType = "movie"

	// A single episode in a TV series
	ContentTypeEpisode ContentType = "episode"

	// A complete season container (used in OTT libraries)
	ContentTypeSeason ContentType = "season"

	// Movie or series trailer
	ContentTypeTrailer ContentType = "trailer"

	// Promotional/teaser videos (OTT promos, announcements)
	ContentTypePromo ContentType = "promo"

	// A playlist of videos (YouTube playlist, Twitch collections)
	ContentTypePlaylist ContentType = "playlist"

	// Channel/uploader/creator identity (YouTube/Twitch channel object)
	ContentTypeChannel ContentType = "channel"

	// Category/genre tags (e.g., Gaming, Sports, Movies)
	ContentTypeCategory ContentType = "category"
)


// ContentItem is the normalized content representation used in the service.
// ContentItem represents a single piece of aggregated content
// from any provider (YouTube, Twitch, OTT, etc.).
type ContentItem struct {
	ID           string      `json:"id"`                    // Unique content ID (provider video/clip/episode ID)
	Provider     string      `json:"provider"`              // Name of provider: youtube, twitch, hotstar, etc.
	ContentType  ContentType `json:"contentType"`           // Normalized content type (video, live, vod, clip, movie, etc.)

	Title        string      `json:"title,omitempty"`       // Title of the content
	Description  string      `json:"description,omitempty"` // Description or summary
	PublishedAt  PublishedAt  `json:"publishedAt,omitempty"` // When the content was published (nil if unknown)
	ThumbnailURL string      `json:"thumbnailUrl,omitempty"`// Best-selected thumbnail for TV UI
	ContentURL   string      `json:"contentUrl,omitempty"`  // Direct link to open/watch the content
	ViewCount    ViewCount       `json:"viewCount,omitempty"`   // Total view count if provided by source

	ChannelID    string      `json:"channelId,omitempty"`   // Provider channel/uploader ID
	Channel      string      `json:"channelTitle,omitempty"`// Channel/uploader display name
	ChannelURL   string      `json:"channelUrl,omitempty"`  // Link to channel/uploader page

	Extras       *Extras      `json:"extras,omitempty"`     // Optional provider-specific metadata (season/episode/tags/etc.)
}

type Extras struct {
	DurationSeconds *Duration   `json:"durationSeconds,omitempty"` // Optional duration override when provider gives extra precision
	Language        *string  `json:"language,omitempty"`        // Content language (en, hi, te, etc.)
	Rating          *string  `json:"rating,omitempty"`          // Content rating (U/A, PG-13, 18+ etc.)
	Tags            []string `json:"tags,omitempty"`            // Provider tags/keywords (YouTube tags, genres)

	SeasonNumber    *int     `json:"seasonNumber,omitempty"`    // Applicable for OTT episodes (Sx)
	EpisodeNumber   *int     `json:"episodeNumber,omitempty"`   // Applicable for OTT episodes (Ex)

	IsLive          *bool    `json:"isLive,omitempty"`          // Provider live flag (true for livestream content)
	IsPremium       *bool    `json:"isPremium,omitempty"`       // OTT premium-status (requires subscription)
}

type ViewCount struct {
	Value   int64  `json:"value,omitempty"`   // numeric value
	Display string `json:"display,omitempty"` // formatted string: "2.3K"
}

type PublishedAt struct {
	Timestamp time.Time `json:"timestamp,omitempty"`
	Relative  string    `json:"relative,omitempty"`
}

type Duration struct {
    Seconds int64  `json:"seconds,omitempty"`
    Display string `json:"display,omitempty"` // "1:30"
}