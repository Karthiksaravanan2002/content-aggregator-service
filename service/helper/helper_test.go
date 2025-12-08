package helper_test

import (
	"net/http"
	"testing"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/helper"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"

	"github.com/stretchr/testify/assert"
)

func Test_SelectRespStatusCode(t *testing.T) {

	tests := []struct {
		name     string
		resp     *models.AggregateResponse
		expected int
	}{
		{
			name: "ALL SUCCESS should return 200",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data:          map[string][]*models.ContentItem{"trending": {}},
						FeatureErrors: map[string]errors.AppError{},
					},
				},
			},
			expected: http.StatusOK,
		},

		{
			name: "MIXED CASE should return data + errors with code 207",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data: map[string][]*models.ContentItem{
							"trending": {{ID: "yt1"}},
						},
						FeatureErrors: map[string]errors.AppError{
							"mylist": errors.BadRequest(newErr("bad req")),
						},
					},
				},
			},
			expected: http.StatusMultiStatus,
		},

		{
			name: "ALL FAILED should return priority error",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						Data: map[string][]*models.ContentItem{},
						FeatureErrors: map[string]errors.AppError{
							"trending": errors.InternalError(newErr("internal")),
						},
					},
				},
			},
			expected: http.StatusInternalServerError,
		},

		{
			name: "MULTIPLE ERRORS should return lowest priority (400 < 500)",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						FeatureErrors: map[string]errors.AppError{
							"trending": errors.BadRequest(newErr("bad req")),
							"mylist":   errors.InternalError(newErr("internal")),
						},
					},
				},
			},
			expected: http.StatusBadRequest, // priority: 400 < 500
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.SelectRespStatusCode(tt.resp)
			assert.Equal(t, tt.expected, got)
		})
	}
}

func Test_SelectPriorityError(t *testing.T) {

	tests := []struct {
		name     string
		resp     *models.AggregateResponse
		expected int
	}{
		{
			name: "choose BAD_REQUEST over INTERNAL_ERROR",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"twitch": {
						FeatureErrors: map[string]errors.AppError{
							"f1": errors.InternalError(newErr("internal")),
							"f2": errors.BadRequest(newErr("bad req")),
						},
					},
				},
			},
			expected: http.StatusBadRequest,
		},

		{
			name: "choose BAD_GATEWAY if no better errors",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"yt": {
						FeatureErrors: map[string]errors.AppError{
							"f1": errors.BadGateway(newErr("upstream down")),
						},
					},
				},
			},
			expected: http.StatusBadGateway,
		},

		{
			name: "NO ERRORS return 500",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"none": {FeatureErrors: map[string]errors.AppError{}},
				},
			},
			expected: http.StatusInternalServerError,
		},

		{
			name: "multiple providers returns best overall priority",
			resp: &models.AggregateResponse{
				Providers: map[string]*models.ProviderResponse{
					"youtube": {
						FeatureErrors: map[string]errors.AppError{
							"a": errors.InternalError(newErr("internal")),
						},
					},
					"twitch": {
						FeatureErrors: map[string]errors.AppError{
							"b": errors.BadRequest(newErr("bad req")),
						},
					},
				},
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := helper.SelectPriorityError(tt.resp)
			assert.Equal(t, tt.expected, got)
		})
	}
}

// simple local error wrapper for tests
func newErr(msg string) error { return &simpleErr{msg} }

type simpleErr struct{ s string }

func (e *simpleErr) Error() string { return e.s }
