package repository

import (
	"context"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
)

const (
	ARTICLES_REPO_VERSION1 = "version1"
	LevelReadCommitted     = 1
	LevelReadUncommitted   = 2
	LevelWriteCommitted    = 3
	LevelRepeatableRead    = 4
	LevelSnapshot          = 5
	LevelSerializable      = 6
	LevelLinearizable      = 7
	LevelDefault           = 0
)

// Todo to reuse the code of transactions instead of doing it always
type ArticlesRepository interface {
	InsertArticle(ctx context.Context, article *database.Article) (int64, error)
	GetArticles(ctx context.Context) ([]*database.Article, error)
	GetArticle(ctx context.Context, query string, args ...interface{}) (*database.Article, error)
	UpdateArticle(ctx context.Context, article *database.Article) (int64, error)
}

type ArticlesRepositoryImpl struct {
	client databases.SqlClient
}

func NewArticlesRepositoryFactory(version string, client databases.SqlClient) ArticlesRepository {
	switch version {
	case ARTICLES_REPO_VERSION1:
		return &ArticlesRepositoryImpl{client: client}
	default:
		return &ArticlesRepositoryImpl{client: client}
	}
}

func (repo *ArticlesRepositoryImpl) InsertArticle(ctx context.Context, article *database.Article) (int64, error) {
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

func (repo *ArticlesRepositoryImpl) GetArticles(ctx context.Context) ([]*database.Article, error) {
	results, err := repo.client.QueryWithContext(ctx, helpers.GetAllArticles)
	if err != nil {
		return nil, err
	}
	finalResult := []*database.Article{}
	for results.Next() {
		article := database.Article{}
		err := results.Scan(&article.ArtID, &article.Name, &article.AvilableStock, &article.SoldStock)
		if err != nil {
			return nil, err
		}
		finalResult = append(finalResult, &article)
	}
	return finalResult, nil
}

func (repo *ArticlesRepositoryImpl) GetArticle(ctx context.Context, query string, args ...interface{}) (*database.Article, error) {
	result := repo.client.QueryRowWithContext(ctx, query, args...)
	article := database.Article{}
	err := result.Scan(&article.ArtID, &article.Name, &article.AvilableStock, &article.SoldStock)
	if err != nil {
		return nil, err
	}
	return &article, nil
}

func (repo *ArticlesRepositoryImpl) UpdateArticle(ctx context.Context, article *database.Article) (int64, error) {
	options := repo.client.BuildOptions(false, repo.client.GetIsolationLevel(1))
	tx, err := repo.client.BeginWithContext(ctx, &options)
	if err != nil {
		return -1, err
	}
	stmt, err := repo.client.PrepareWithContext(ctx, tx, helpers.UpdateArticleByID)
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
