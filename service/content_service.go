package service

import (
	"context"
	"time"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"go.uber.org/zap"
)

type ContentService struct {
	factory *ProviderFactory
	logger  *zap.Logger
}

func NewContentService(factory *ProviderFactory, logger *zap.Logger, timeout time.Duration) *ContentService {
	return &ContentService{
		factory: factory,
		logger:  logger.With(zap.String("component", "content-service")),
	}
}

func (s *ContentService) HandleRequest(ctx context.Context,req models.AggregateRequest) (models.AggregateResponse, error) {

	result := models.AggregateResponse{Providers: make(map[string]models.ProviderResponse)}

	for _, p := range req.Providers {

		provider := s.factory.GetProvider(p.Provider)

		  pResp := models.ProviderResponse{
            Data:          make(map[string][]models.ContentItem),
            FeatureErrors: make(map[string]*errors.AppError),
        }
		if provider == nil {
			s.logger.Warn("unsupported provider", zap.String("provider", p.Provider))
      err := errors.Unsupported("provider " + p.Provider)

            pResp.FeatureErrors["_provider"] = err
						result.Providers[p.Provider] = pResp

			continue


		}

		for _, featureName := range p.Functionality {

			items, err := provider.FetchFeature(ctx, p, featureName)

			if err != nil {
				s.logger.Error("feature fetch failed",zap.String("provider", p.Provider),zap.String("feature", featureName),zap.Error(err),)

				// add feature-specific error
				pResp.FeatureErrors[featureName] = errors.ProviderFailure(p.Provider,err)

				continue
			}

			// success → store items
			pResp.Data[featureName] = items
		}

		// 3 — store provider response
		result.Providers[p.Provider] = pResp
	}

	return result, nil
}
