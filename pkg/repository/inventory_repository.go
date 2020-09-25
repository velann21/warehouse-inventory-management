package repository

import (
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
)

const (
	INVENTORY_REPO_VERSION1 = "version1"
)
type Repository interface {
	InsertArticle(ctx context.Context, article requests.AddArticles)error
}

type InventoryRepository struct {
	client databases.SqlClient
}

func NewInventoryRepositoryFactory(version string, client databases.SqlClient)Repository{
	switch version {
	case INVENTORY_REPO_VERSION1:
		return &InventoryRepository{client:client}
	default:
		return &InventoryRepository{client:client}
	}
}

func (repo *InventoryRepository) InsertArticle(ctx context.Context, article requests.AddArticles)error{
	logrus.Info("Starting AddArticles repo")
	tx, err := repo.client.Begin()
	if err != nil{
		return err
	}

	fmt.Println(tx)
	return nil
}


