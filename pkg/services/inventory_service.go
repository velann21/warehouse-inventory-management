package services

import (
	"context"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
)
const (
	SERVICE_VERSION1 = "Version1"
	)

type Service interface {
	AddArticles(ctx context.Context, articles requests.AddArticles)error
}

type InventoryService struct {
	repo repository.Repository
}

func NewInventoryServiceFactory(version string, repo repository.Repository)Service{
	switch version {
	case SERVICE_VERSION1:
		return &InventoryService{repo:repo}
	default:
		return &InventoryService{repo:repo}
	}
}

func (inventoryService *InventoryService) AddArticles(ctx context.Context, articles requests.AddArticles)error{
	err := inventoryService.repo.InsertArticle(ctx, articles)
	if err != nil{
		return err
	}
	return nil
}
