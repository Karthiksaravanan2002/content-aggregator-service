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

func (s *ContentService) HandleRequest(ctx context.Context, req models.AggregateRequest) (*models.AggregateResponse, error) {
	
	
	result := &models.AggregateResponse{
		Providers: make(map[string]models.ProviderResponse),
	}

	for _, p := range req.Providers {

		provider := s.factory.GetProvider(p.Provider)
		
		if provider == nil {
			result.Providers[p.Provider] = models.ProviderResponse{ 
				FeatureErrors: map[string]*errors.AppError{"_provider": errors.NewBadRequest("unsupported provider: "+p.Provider,map[string]interface{}{"provider": p.Provider},
			   ),
				},
			}
			continue
		}

		featureData := make(map[string][]models.ContentItem)
		featureErrs := make(map[string]*errors.AppError)

		for _, featName := range p.Functionality {

			feat := FeatureRegistry[featName]
			if feat == nil {
				featureErrs[featName] = errors.NewBadRequest("unsupported feature: "+featName,map[string]interface{}{"feature": featName},)
				continue
			}

			raw, err := provider.FetchFeatureRaw(ctx, p, featName)

			
			if err != nil {
				featureErrs[featName] = errors.NewProviderError(p.Provider,"Failed to fetch provider",map[string]interface{}{"feature": featName, "error": err.Error()},)
				continue
			}

			items, ferr := feat.Apply(ctx, raw)
				if ferr != nil {featureErrs[featName] = errors.NewFeatureError(featName,"feature execution failed",map[string]interface{}{"error": ferr.Error()},)
				continue
			}

		featureData[featName] = items
		}

		result.Providers[p.Provider] = models.ProviderResponse{
			Data:  featureData,
			FeatureErrors: featureErrs,
		}
	}

	return result, nil
}
