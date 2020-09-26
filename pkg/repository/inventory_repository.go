package repository

import (
	"context"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
)

const (
	INVENTORY_REPO_VERSION1 = "version1"
	LevelReadCommitted      = 1
	LevelReadUncommitted    = 2
	LevelWriteCommitted     = 3
	LevelRepeatableRead     = 4
	LevelSnapshot           = 5
	LevelSerializable       = 6
	LevelLinearizable       = 7
	LevelDefault            = 0
)

type InventoryRepository interface {
	InsertArticle(ctx context.Context, article *database.Article) (int64, error)
	GetArticles(ctx context.Context) ([]database.Article, error)
	GetArticle(ctx context.Context, query string, args ...interface{}) (*database.Article, error)
	UpdateArticle(ctx context.Context, article *database.Article) (int64, error)
}

type InventoryRepositoryImpl struct {
	client databases.SqlClient
}

func NewInventoryRepositoryFactory(version string, client databases.SqlClient) InventoryRepository {
	switch version {
	case INVENTORY_REPO_VERSION1:
		return &InventoryRepositoryImpl{client: client}
	default:
		return &InventoryRepositoryImpl{client: client}
	}
}

func (repo *InventoryRepositoryImpl) InsertArticle(ctx context.Context, article *database.Article) (int64, error) {
	options := repo.client.BuildOptions(false, repo.client.GetIsolationLevel(1))
	tx, err := repo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}

	stmt, err := repo.client.PrepareWithContext(ctx, tx, helpers.CreateArticle)
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	res, err := repo.client.ExecWithContext(ctx, stmt, article.GetName(), article.GetAvailableStock(), article.GetSoldStock())
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	id, err := repo.client.LastInsertedID(res)
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	err = repo.client.Commit(tx)
	if err != nil {
		err = repo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}

func (repo *InventoryRepositoryImpl) GetArticles(ctx context.Context) ([]database.Article, error) {
	results, err := repo.client.QueryWithContext(ctx, helpers.CreateArticle)
	if err != nil {
		return nil, err
	}
	finalResult := []database.Article{}
	for results.Next() {
		article := database.Article{}
		err := results.Scan(article.SetArtID, article.SetName, article.SetAvailableStock, article.SetSoldStock)
		if err != nil {
			return nil, err
		}
		finalResult = append(finalResult, article)
	}
	return finalResult, nil
}

func (repo *InventoryRepositoryImpl) GetArticle(ctx context.Context, query string, args ...interface{}) (*database.Article, error) {
	result := repo.client.QueryRowWithContext(ctx, query, args...)
	article := database.Article{}
	err := result.Scan(&article.ArtID, &article.Name, &article.AvilableStock, &article.SoldStock)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (repo *InventoryRepositoryImpl) UpdateArticle(ctx context.Context, article *database.Article) (int64, error) {
	options := repo.client.BuildOptions(false, repo.client.GetIsolationLevel(1))
	tx, err := repo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := repo.client.PrepareWithContext(ctx, tx, helpers.UpdateArticle)
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	res, err := repo.client.ExecWithContext(ctx, stmt, article.GetName(), article.GetAvailableStock(), article.GetSoldStock(), article.GetArtID())
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	id, err := repo.client.LastInsertedID(res)
	if err != nil {
		_ = repo.client.RollBack(tx)
		return -1, err
	}
	err = repo.client.Commit(tx)
	if err != nil {
		err = repo.client.RollBack(tx)
		if err != nil {
			return -1, err
		}
		return -1, err
	}
	return id, nil
}
