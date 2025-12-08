package service

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"

	"dev.azure.com/daimler-mic/content-aggregator/service/cache"
	"dev.azure.com/daimler-mic/content-aggregator/service/errors"
	"dev.azure.com/daimler-mic/content-aggregator/service/models"
	"dev.azure.com/daimler-mic/content-aggregator/service/props"
	"dev.azure.com/daimler-mic/content-aggregator/service/providers"
)

type dummyProvider struct{}

func (d *dummyProvider) FetchFeature(ctx context.Context, req *models.ProviderRequest, feature string) ([]*models.ContentItem, errors.AppError) {
	return nil, nil
}

func (d *dummyProvider) GetFeature(name string) providers.FeatureStrategy {
	return nil
}

func TestProviderFactory(t *testing.T) {

	logger := zap.NewNop()

	noCache := cache.NewNoopCache()
	cfg := props.ProvidersConfig{}
	cacheCfg := props.CacheConfig{}

	tests := []struct {
		name            string
		setupFactory    func() ProviderFactory
		providerName    string
		expectNil       bool
		expectDecorator bool
	}{
		{
			name: "unknown provider should return nil",
			setupFactory: func() ProviderFactory {
				return NewProviderFactory(cfg, cacheCfg, noCache, logger)
			},
			providerName:    "does-not-exist",
			expectNil:       true,
			expectDecorator: false,
		},
		{
			name: "registered provider returns wrapped provider",
			setupFactory: func() ProviderFactory {

				f := &providerFactory{
					props:    cfg,
					cacheCfg: cacheCfg,
					cache:    noCache,
					logger:   logger,
					makers:   make(map[string]ProviderMaker),
				}

				f.Register("test", func(l *zap.Logger) providers.ProviderStrategy {
					return &dummyProvider{}
				})

				return f
			},
			providerName:    "test",
			expectNil:       false,
			expectDecorator: true,
		},
		{
			name: "case-insensitive provider lookup works",
			setupFactory: func() ProviderFactory {

				f := &providerFactory{
					props:    cfg,
					cacheCfg: cacheCfg,
					cache:    noCache,
					logger:   logger,
					makers:   make(map[string]ProviderMaker),
				}

				f.Register("YouTube", func(l *zap.Logger) providers.ProviderStrategy {
					return &dummyProvider{}
				})

				return f
			},
			providerName:    "youtube",
			expectNil:       false,
			expectDecorator: true,
		},
	}

	for _, tt := range tests {

		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			factory := tt.setupFactory()
			provider := factory.GetProvider(tt.providerName)

			// If expecting nil
			if tt.expectNil {
				assert.Nil(t, provider)
				return
			}

			// Else expecting a decorated provider
			require.NotNil(t, provider)

			if tt.expectDecorator {
				_, ok := provider.(*cache.CacheDecorator)
				assert.True(t, ok, "provider must be wrapped inside CacheDecorator")
			}
		})
	}
}
