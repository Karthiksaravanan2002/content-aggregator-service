package models

import "time"

// AggregateRequest represents the incoming request from the client.
type AggregateRequest struct {
	Providers []ProviderRequest `json:"providers"`
}

// ProviderRequest represents a single provider invocation request.
type ProviderRequest struct {
	Provider      string   `json:"provider"`
	APIKey        string   `json:"apiKey,omitempty"`
	Functionality []string `json:"functionality,omitempty"`
}

// ContentItem is the normalized content representation used in the service.
type ContentItem struct {
	ID          string            `json:"id"`
	Provider    string            `json:"provider"`
	Title       string            `json:"title,omitempty"`
	Description string            `json:"description,omitempty"`
	PublishedAt *time.Time        `json:"publishedAt,omitempty"`
	ChannelID   string            `json:"channelId,omitempty"`
	Channel     string            `json:"channelTitle,omitempty"`
	Thumbnail   map[string]string `json:"thumbnail,omitempty"`
	ViewCount   int64             `json:"viewCount,omitempty"`
	Extra       map[string]string `json:"extra,omitempty"`
}
