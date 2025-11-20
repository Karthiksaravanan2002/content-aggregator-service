package models

import "dev.azure.com/daimler-mic/content-aggregator/service/errors"


type AggregateResponse struct {
	Providers map[string]*ProviderResponse `json:"providers"`
}

type ProviderResponse struct {
    Data          map[string][]*ContentItem   `json:"data,omitempty"`
    FeatureErrors map[string]errors.AppError `json:"featureErrors,omitempty"`
}
