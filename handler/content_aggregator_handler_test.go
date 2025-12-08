package handler_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"dev.azure.com/daimler-mic/content-aggregator/handler"
	"dev.azure.com/daimler-mic/content-aggregator/service"
	appErr "dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

type mockContentService struct {
	resp *models.AggregateResponse
}

func (m *mockContentService) Aggregate(ctx context.Context, req *models.AggregateRequest) *models.AggregateResponse {
	return m.resp
}

func Test_ContentHandler_ServeHTTP(t *testing.T) {

	tests := []struct {
		name           string
		service        *mockContentService
		props          *props.Config
		requestMethod  string
		requestBody    string
		expectedStatus int
		expectErrorMsg string
		expectedData   map[string]map[string]string // provider → feature → expectedTitle
	}{
		{
			name:           "invalid http method",
			service:        &mockContentService{},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodGet,
			requestBody:    "",
			expectedStatus: http.StatusMethodNotAllowed,
		},

		{
			name:           "invalid JSON",
			service:        &mockContentService{},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodPost,
			requestBody:    "{invalid",
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "invalid request body",
		},

		{
			name:           "no providers supplied",
			service:        &mockContentService{},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodPost,
			requestBody:    `{"providers":[]}`,
			expectedStatus: http.StatusBadRequest,
			expectErrorMsg: "at least one provider must be supplied",
		},

		{
			name: "success response",
			service: &mockContentService{
				resp: &models.AggregateResponse{
					Providers: map[string]*models.ProviderResponse{
						"youtube": {
							Data: map[string][]*models.ContentItem{
								"trending": {
									{ID: "yt1", Title: "Video 1"},
								},
							},
						},
					},
				},
			},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodPost,
			requestBody:    `{"providers":[{"provider":"youtube","functionality":["trending"]}]}`,
			expectedStatus: http.StatusOK,
			expectedData: map[string]map[string]string{
				"youtube": {"trending": "Video 1"},
			},
		},

		{
			name: "multi-status response",
			service: &mockContentService{
				resp: &models.AggregateResponse{
					Providers: map[string]*models.ProviderResponse{
						"youtube": {
							Data: map[string][]*models.ContentItem{
								"trending": {{ID: "yt1", Title: "Video 1"}},
							},
							FeatureErrors: map[string]appErr.AppError{
								"mylist": appErr.BadGateway(assert.AnError),
							},
						},
					},
				},
			},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodPost,
			requestBody:    `{"providers":[{"provider":"youtube","functionality":["trending","mylist"]}]}`,
			expectedStatus: http.StatusMultiStatus,
		},

		{
			name: "all failed -> priority error returned",
			service: &mockContentService{
				resp: &models.AggregateResponse{
					Providers: map[string]*models.ProviderResponse{
						"youtube": {
							FeatureErrors: map[string]appErr.AppError{
								"trending": appErr.BadRequest(errors.New("invalid params")),
							},
						},
						"twitch": {
							FeatureErrors: map[string]appErr.AppError{
								"trending": appErr.BadGateway(errors.New("provider timeout")),
							},
						},
					},
				},
			},
			props:          &props.Config{Server: props.ServerConfig{Timeout: 1}},
			requestMethod:  http.MethodPost,
			requestBody:    `{"providers":[{"provider":"youtube","functionality":["trending"]}]}`,
			expectedStatus: http.StatusBadRequest, // 400 wins over 502
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			logger, _ := zap.NewDevelopment()

			// Build handler with mock service
			h := handler.NewContentHandler((service.ContentService)(tt.service),
				logger,
				tt.props,
			)

			// Prepare request
			req := httptest.NewRequest(tt.requestMethod, "/content/aggregate", bytes.NewBufferString(tt.requestBody))
			rr := httptest.NewRecorder()

			// Execute handler
			h.ServeHTTP(rr, req)

			// Validate status code
			assert.Equal(t, tt.expectedStatus, rr.Code, "unexpected status code")

			// Validate error message (plain-text error responses)
			if tt.expectErrorMsg != "" {
				assert.Contains(t, rr.Body.String(), tt.expectErrorMsg)
				return
			}

			// Validate JSON response on success/multi-status
			if tt.expectedData != nil {
				var resp models.AggregateResponse
				err := json.Unmarshal(rr.Body.Bytes(), &resp)
				assert.NoError(t, err)

				for prov, feats := range tt.expectedData {
					for feat, expectedTitle := range feats {
						items := resp.Providers[prov].Data[feat]
						assert.NotNil(t, items)
						assert.Equal(t, expectedTitle, items[0].Title)
					}
				}
			}
		})
	}
}
