package service

import (
	"context"
	"time"

	stdErr "errors"

	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"go.uber.org/zap"
)

type ContentService interface{
	Aggregate(ctx context.Context,req *models.AggregateRequest) (*models.AggregateResponse)
}

type contentService struct {
	factory ProviderFactory
	logger  *zap.Logger
} 

func (s *contentService) Aggregate(ctx context.Context,req *models.AggregateRequest) (*models.AggregateResponse) {

	result := &models.AggregateResponse{Providers: make(map[string]*models.ProviderResponse)}

	for _, p := range req.Providers {

		provider := s.factory.GetProvider(p.Provider)

		  pResp := &models.ProviderResponse{
            Data:          make(map[string][]*models.ContentItem),
            FeatureErrors: make(map[string]errors.AppError),
        }
		if provider == nil {
			s.logger.Warn("unsupported provider", zap.String("provider", p.Provider))
			

            pResp.FeatureErrors["Default"] = errors.BadRequest(stdErr.New("Provider not supported"))
						result.Providers[p.Provider] = pResp

			continue
		}

		if p.Functionality == nil{
			      pResp.FeatureErrors["Default"] = errors.BadRequest(stdErr.New("No Features Present"))
						result.Providers[p.Provider] = pResp
		}

		for _, featureName := range p.Functionality {

			items, err := provider.FetchFeature(ctx, &p, featureName)
			if err != nil {
				s.logger.Error("feature fetch failed",zap.String("provider", p.Provider),zap.String("feature", featureName),zap.Error(err),)
				pResp.FeatureErrors[featureName] = err
				continue
			}
			pResp.Data[featureName] = items
		}
		result.Providers[p.Provider] = pResp
	}

	return result
}


func NewContentService(factory ProviderFactory, logger *zap.Logger, timeout time.Duration) ContentService {
	return &contentService{
		factory: factory,
		logger:  logger.With(zap.String("component", "content-service")),
	}
}