package services

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	sqlMigration"github.com/velann21/warehouse-inventory-management/pkg/migration_scripts"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
	articlesModel "github.com/velann21/warehouse-inventory-management/pkg/models/internals"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
	"strconv"
)

const (
	ARTICLES_SERVICE_VERSION1 = "Version1"
)

type ArticlesService interface {
	AddArticles(ctx context.Context, articles *requests.AddArticles) ([]articlesModel.SuccessfullyAddedArticle, []articlesModel.FailedArticle)
	AddArticlesFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error
	ListArticles(ctx context.Context) ([]*database.Article, error)
	SqlMigration(ctx context.Context, req *requests.SqlMigrationRequest)error

	// Private functions
	assignTask(ctx context.Context, articleStreams chan *requests.Article, waitChannel chan bool)
	addArticleJob(ctx context.Context, articleData *requests.Article) (int64, error)
	streamArticlesData(decoder *json.Decoder) chan *requests.Article
	createNewArticle(ctx context.Context, value *requests.Article) (int64, error)
	addArticleAsSuccess(article *requests.Article, suc []articlesModel.SuccessfullyAddedArticle, id int64)
	addArticleAsFailure(article *requests.Article, errorsArr []articlesModel.FailedArticle, err error)
}

type ArticlesServiceImpl struct {
	repo repository.ArticlesRepository
}

func NewInventoryServiceFactory(version string, repo repository.ArticlesRepository) ArticlesService {
	switch version {
	case ARTICLES_SERVICE_VERSION1:
		return &ArticlesServiceImpl{repo: repo}
	default:
		return &ArticlesServiceImpl{repo: repo}
	}
}

func (articlesService *ArticlesServiceImpl) AddArticles(ctx context.Context, articles *requests.AddArticles) ([]articlesModel.SuccessfullyAddedArticle, []articlesModel.FailedArticle) {
	succedArticles := make([]articlesModel.SuccessfullyAddedArticle, 0)
	failedArticles := make([]articlesModel.FailedArticle, 0)
	for _, value := range articles.Articles {
		id, err := articlesService.addArticleJob(ctx, &value)
		if err != nil {
			logrus.WithError(err).Error("Something went wrong while addArticleJob for" + value.Name)
			articlesService.addArticleAsFailure(&value, failedArticles, err)
			continue
		}
		articlesService.addArticleAsSuccess(&value, succedArticles, id)
	}
	return succedArticles, failedArticles
}

func (articlesService *ArticlesServiceImpl) AddArticlesFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error {
	if decoder == nil {
		return nil
	}
	articleStreams := articlesService.streamArticlesData(decoder)
	articlesService.assignTask(ctx, articleStreams, waitChannel)
	return nil
}

func (articlesService *ArticlesServiceImpl) assignTask(ctx context.Context, articleStreams chan *requests.Article, waitChannel chan bool) {
	go func() {
		succedArticles := make([]articlesModel.SuccessfullyAddedArticle, 0)
		failedArticles := make([]articlesModel.FailedArticle, 0)
		for articleData := range articleStreams {
			id, err := articlesService.addArticleJob(ctx, articleData)
			if err != nil {
				logrus.WithError(err).Error("Something went wrong while addArticleJob for" + articleData.Name)
				articlesService.addArticleAsFailure(articleData, failedArticles, err)
				continue
			}
			articlesService.addArticleAsSuccess(articleData, succedArticles, id)
		}
		waitChannel <- true
	}()
}

func (articlesService *ArticlesServiceImpl) streamArticlesData(decoder *json.Decoder) chan *requests.Article {
	streamChannel := make(chan *requests.Article, 15000)
	go func() {
		_, _ = decoder.Token()
		for decoder.More() {
			articleObj := requests.Article{}
			err := decoder.Decode(&articleObj)
			if err != nil {
				logrus.WithError(err).Error("Something wrong while decode inside streamArticlesData()")
				continue
			}
			streamChannel <- &articleObj
		}
		close(streamChannel)
	}()
	return streamChannel
}

func (articlesService *ArticlesServiceImpl) addArticleJob(ctx context.Context, articleData *requests.Article) (int64, error) {
	article, err := articlesService.getArticle(ctx, articleData)
	if err != nil {
		if err.Error() == helpers.SQLRowNotFound {
			id, err := articlesService.createNewArticle(ctx, articleData)
			if err != nil {
				return id, err
			}
			return id, nil

		}
		return -1, err
	}
	if article != nil {
		convertedData, err := strconv.Atoi(articleData.Stock)
		if err != nil {

		}
		article.AvilableStock += convertedData
		updatedID, err := articlesService.repo.UpdateArticle(ctx, article)
		if err != nil {
			return updatedID, err
		}
		return updatedID, nil
	}
	return -1, nil
}

func (articlesService *ArticlesServiceImpl) createNewArticle(ctx context.Context, value *requests.Article) (int64, error) {
	stock, err := strconv.Atoi(value.Stock)
	newArticle := database.NewArticle(0, stock, value.Name)
	id, err := articlesService.repo.InsertArticle(ctx, newArticle)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (articlesService *ArticlesServiceImpl) getArticle(ctx context.Context, value *requests.Article) (*database.Article, error) {
	return articlesService.repo.GetArticle(ctx, helpers.GetArticleByName, value.Name)
}

func (articlesService *ArticlesServiceImpl) addArticleAsSuccess(article *requests.Article, suc []articlesModel.SuccessfullyAddedArticle, id int64) {
	success := articlesModel.SuccessfullyAddedArticle{}
	success.Name = article.Name
	success.Total = article.Stock
	success.Endpoint = "/api/v1/inventory/articles/" + strconv.Itoa(int(int64(id)))
	suc = append(suc, success)
}

func (articlesService *ArticlesServiceImpl) addArticleAsFailure(article *requests.Article, errorsArr []articlesModel.FailedArticle, err error) {
	error := articlesModel.FailedArticle{}
	error.Name = article.Name
	error.Reason = err.Error()
	errorsArr = append(errorsArr, error)
}

func (articlesService *ArticlesServiceImpl) ListArticles(ctx context.Context) ([]*database.Article, error) {
	articles, err := articlesService.repo.GetArticles(ctx)
	if err != nil {
		return nil, err
	}
	return articles, nil
}


func (articlesService *ArticlesServiceImpl) SqlMigration(ctx context.Context, req *requests.SqlMigrationRequest)error{
	if req.Upcount > 0{
		err := sqlMigration.MigrateDb(uint(req.Upcount))
		if err != nil{
			return err
		}
	}else if req.Downcount < 0{
		err := sqlMigration.MigrateDb(uint(req.Downcount))
		if err != nil{
			return err
		}
	}
	return nil
}
