package service

import (
	"dev.azure.com/daimler-mic/content-aggregator/service/features"
	"dev.azure.com/daimler-mic/content-aggregator/service/features/continuewatching"
	"dev.azure.com/daimler-mic/content-aggregator/service/features/mylist"
	"dev.azure.com/daimler-mic/content-aggregator/service/features/trending"
)

// FeatureRegistry maps feature names to strategy implementations.
var FeatureRegistry = map[string]features.FeatureStrategy{
	"trending":          trending.NewTrendingFeature(),
	"mylist":            mylist.NewMyListFeature(),
	"continue_watching": continuewatching.NewContinueWatchingFeature(),
}
