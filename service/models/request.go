package models

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
