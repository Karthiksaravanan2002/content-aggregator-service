package mylist

import (
	"context"

	"dev.azure.com/daimler-mic/content-aggregator/service/models"
)

type MyListFeature struct{}

func NewMyListFeature() *MyListFeature { return &MyListFeature{} }
func (MyListFeature) Name() string     { return "mylist" }

func (MyListFeature) Apply(ctx context.Context, raw map[string][]models.ContentItem) ([]models.ContentItem, error) {
	return []models.ContentItem{}, nil
}
