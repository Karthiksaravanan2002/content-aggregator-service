package service

import (
	"context"
	"errors"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"go.uber.org/zap"
)

// ContentService aggregates provider content & applies features.
type ContentService struct {
	factory  *ProviderFactory
	logger   *zap.Logger
	timeout  time.Duration
	cacheTTL time.Duration
}

func (s *ContentService) HandleRequest(ctx context.Context, req models.AggregateRequest) (map[string]interface{}, error) {

	if len(req.Providers) == 0 {
		return nil, errors.New("no providers supplied")
	}

	s.logger.Info("content aggregation started",
		zap.Int("providers", len(req.Providers)),
	)

	response := make(map[string]interface{})

	for _, p := range req.Providers {
		provider := s.factory.GetProvider(p.Provider)
		if provider == nil {
			response[p.Provider] = map[string]string{
				"error": "unsupported provider",
			}
			continue
		}

		// Timeout per provider call
		pctx, cancel := context.WithTimeout(ctx, s.timeout)
		defer cancel()

		rawData, err := provider.FetchContent(pctx, p)
		if err != nil {
			s.logger.Error("provider fetch failed", zap.String("provider", p.Provider), zap.Error(err))
			response[p.Provider] = map[string]string{
				"error": err.Error(),
			}
			continue
		}

		// rawData is map[string][]ContentItem
		featureOutput := make(map[string][]models.ContentItem)

		for _, featureName := range p.Functionality {

			feature, ok := FeatureRegistry[featureName]
			if !ok {
				featureOutput[featureName] = []models.ContentItem{
					//{ID: "error", Raw: map[string]interface{}{"error": "unsupported feature"}},
				}
				continue
			}

			// If provider didn't provide that feature key, skip or err
			items, ok := rawData[featureName]
			if !ok {
				featureOutput[featureName] = []models.ContentItem{
					//{ID: "error", Raw: map[string]interface{}{"error": "provider did not return feature"}},
				}
				continue
			}

			// Apply feature transformation
			out := feature.Execute(ctx, items, p)
			featureOutput[featureName] = out
		}

		response[p.Provider] = featureOutput
	}

	return response, nil
}

func NewContentService(
	factory *ProviderFactory,
	logger *zap.Logger,
	timeout time.Duration,
	cacheTTL time.Duration,
) *ContentService {

	return &ContentService{
		factory:  factory,
		logger:   logger,
		timeout:  timeout,
		cacheTTL: cacheTTL,
	}
}
