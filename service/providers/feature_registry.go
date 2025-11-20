package providers

type FeatureRegistry interface {
	GetFeature(name string) FeatureStrategy
}
