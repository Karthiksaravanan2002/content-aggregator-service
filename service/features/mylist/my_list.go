package mylist

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

// MyListFeature is a placeholder for a user's saved content list.
type MyListFeature struct{}

func NewMyListFeature() *MyListFeature {
	return &MyListFeature{}
}

func (f *MyListFeature) Name() string {
	return "mylist"
}

func (f *MyListFeature) Execute(ctx context.Context, items []models.ContentItem, req models.ProviderRequest) []models.ContentItem {

	return items
}
