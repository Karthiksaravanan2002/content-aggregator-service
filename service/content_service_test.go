package service

import (
	"context"
	stdErr "errors"
	"testing"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
	"github.com/go-playground/assert/v2"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
)

type mockProvider struct {
	responses map[string][]*models.ContentItem
	errors    map[string]errors.AppError
}

func (m *mockProvider) FetchFeature(ctx context.Context, req *models.ProviderRequest, feature string) ([]*models.ContentItem, errors.AppError) {

	if err, ok := m.errors[feature]; ok {
		return nil, err
	}
	return m.responses[feature], nil

}

type mockProviderFactory struct {
	providers map[string]providers.ProviderStrategy
}

func (m *mockProviderFactory) GetProvider(name string) providers.ProviderStrategy {
	return m.providers[name]
}

func TestContentService_Aggregate(t *testing.T) {

	tests := []struct {
		name string

		service ContentService
		request *models.AggregateRequest

		expectedResponse *models.AggregateResponse
	}{
		{
			name: "success response for youtube trending",
			service: &contentService{
				factory: &mockProviderFactory{
					providers: map[string]providers.ProviderStrategy{
						"youtube": &mockProvider{
							responses: map[string][]*models.ContentItem{
								"trending": {
									{ID: "yt1", Provider: "youtube"},
								},
							},
						},
					},
				},
				logger: zap.NewNop(),
			},
			request: &models.AggregateRequest{
				Providers: []models.ProviderRequest{
					{
						Provider:      "youtube",
						Functionality: []string{"trending"},
					},
				},
			},
			expectedResponse: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data: map[string][]*models.ContentItem{
							"trending": {
								{ID: "yt1", Provider: "youtube"},
							},
						},
						FeatureErrors: map[string]errors.AppError{},
					},
				},
			},
		},
		{
			name: "unsupported provider",
			service: &contentService{
				factory: &mockProviderFactory{
					providers: map[string]providers.ProviderStrategy{},
				},
				logger: zap.NewNop(),
			},
			request: &models.AggregateRequest{
				Providers: []models.ProviderRequest{
					{
						Provider:      "unknown",
						Functionality: []string{"trending"},
					},
				},
			},
			expectedResponse: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"unknown": {
						Data: map[string][]*models.ContentItem{},
						FeatureErrors: map[string]errors.AppError{
							"_provider": errors.BadRequest(stdErr.New("Provider not supported")),
						},
					},
				},
			},
		},
		{
			name: "feature fetch error",
			service: &contentService{
				factory: &mockProviderFactory{
					providers: map[string]providers.ProviderStrategy{
						"youtube": &mockProvider{
							errors: map[string]errors.AppError{
								"trending": errors.BadGateway(stdErr.New("provider down")),
							},
						},
					},
				},
				logger: zap.NewNop(),
			},
			request: &models.AggregateRequest{
				Providers: []models.ProviderRequest{
					{
						Provider:      "youtube",
						Functionality: []string{"trending"},
					},
				},
			},
			expectedResponse: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data: map[string][]*models.ContentItem{},
						FeatureErrors: map[string]errors.AppError{
							"trending": errors.BadGateway(stdErr.New("provider down")),
						},
					},
				},
			},
		},
		{
			name: "mixed success and error for multiple features",
			service: &contentService{
				factory: &mockProviderFactory{
					providers: map[string]providers.ProviderStrategy{
						"youtube": &mockProvider{
							responses: map[string][]*models.ContentItem{
								"trending": {
									{ID: "yt1"},
								},
							},
							errors: map[string]errors.AppError{
								"mylist": errors.BadRequest(stdErr.New("invalid token")),
							},
						},
					},
				},
				logger: zap.NewNop(),
			},
			request: &models.AggregateRequest{
				Providers: []models.ProviderRequest{
					{
						Provider:      "youtube",
						Functionality: []string{"trending", "mylist"},
					},
				},
			},
			expectedResponse: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data: map[string][]*models.ContentItem{
							"trending": {
								{ID: "yt1"},
							},
						},
						FeatureErrors: map[string]errors.AppError{
							"mylist": errors.BadRequest(stdErr.New("invalid token")),
						},
					},
				},
			},
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			resp := tt.service.Aggregate(context.Background(), tt.request)

			require.Len(t, resp.Providers, len(tt.expectedResponse.Providers))

			for providerName, expectedProv := range tt.expectedResponse.Providers {

				actualProv, ok := resp.Providers[providerName]
				require.True(t, ok, "missing provider response: %s", providerName)

				// Assert data ONLY if expected data exists
				if len(expectedProv.Data) > 0 {
					assertFeatureData(t, expectedProv.Data, actualProv.Data)
				}

				// Always assert feature errors
				assertFeatureErrors(t, expectedProv.FeatureErrors, actualProv.FeatureErrors)
			}
		})
	}
}

func assertFeatureErrors(t *testing.T, expected map[string]errors.AppError, actual map[string]errors.AppError) {

	require.Len(t, actual, len(expected))

	for feature, expErr := range expected {
		actErr, ok := actual[feature]
		require.True(t, ok, "missing feature error: %s", feature)

		assert.Equal(t, expErr.StatusCode(), actErr.StatusCode())
		assert.Equal(t, expErr.Code(), actErr.Code())
		assert.Equal(t, expErr.Message(), actErr.Message())
	}
}

func assertFeatureData(t *testing.T, expected map[string][]*models.ContentItem, actual map[string][]*models.ContentItem) {

	require.Len(t, actual, len(expected))

	for feature, expItems := range expected {
		actItems, ok := actual[feature]
		require.True(t, ok, "missing data for feature: %s", feature)

		require.Len(t, actItems, len(expItems))
		assert.Equal(t, expItems, actItems)
	}
}
