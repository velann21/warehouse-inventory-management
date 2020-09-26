package services

import (
	"context"
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
	inventoryModel "github.com/velann21/warehouse-inventory-management/pkg/models/internals"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
	"strconv"
)

const (
	SERVICE_VERSION1 = "Version1"
)

type InventoryService interface {
	AddArticles(ctx context.Context, articles *requests.AddArticles) ([]inventoryModel.SuccessfullyAddedArticle, []inventoryModel.FailedArticle)
	AddArticlesFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error
	assignTask(ctx context.Context, articleStreams chan *requests.Article, waitChannel chan bool)
	addArticleJob(ctx context.Context, articleData *requests.Article) (int64, error)
	streamArticlesData(decoder *json.Decoder) chan *requests.Article
	createNewArticle(ctx context.Context, value *requests.Article) (int64, error)
	addArticleAsSuccess(article *requests.Article, suc []inventoryModel.SuccessfullyAddedArticle, id int64)
	addArticleAsFailure(article *requests.Article, errorsArr []inventoryModel.FailedArticle, err error)
}

type InventoryServiceImpl struct {
	repo repository.InventoryRepository
}

func NewInventoryServiceFactory(version string, repo repository.InventoryRepository) InventoryService {
	switch version {
	case SERVICE_VERSION1:
		return &InventoryServiceImpl{repo: repo}
	default:
		return &InventoryServiceImpl{repo: repo}
	}
}

func (inventoryService *InventoryServiceImpl) AddArticles(ctx context.Context, articles *requests.AddArticles) ([]inventoryModel.SuccessfullyAddedArticle, []inventoryModel.FailedArticle) {
	succedArticles := make([]inventoryModel.SuccessfullyAddedArticle, 0)
	failedArticles := make([]inventoryModel.FailedArticle, 0)
	for _, value := range articles.Articles {
		id, err := inventoryService.addArticleJob(ctx, &value)
		if err != nil {
			logrus.WithError(err).Error("Something went wrong while addArticleJob for"+ value.Name)
			inventoryService.addArticleAsFailure(&value, failedArticles, err)
			continue
		}
		inventoryService.addArticleAsSuccess(&value, succedArticles, id)
	}
	return succedArticles, failedArticles
}

func (inventoryService *InventoryServiceImpl) AddArticlesFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error {
	if decoder == nil {
		return nil
	}
	articleStreams := inventoryService.streamArticlesData(decoder)
	inventoryService.assignTask(ctx, articleStreams, waitChannel)
	return nil
}

func (inventoryService *InventoryServiceImpl) assignTask(ctx context.Context, articleStreams chan *requests.Article, waitChannel chan bool) {
	go func() {
		succedArticles := make([]inventoryModel.SuccessfullyAddedArticle, 0)
		failedArticles := make([]inventoryModel.FailedArticle, 0)
		for articleData := range articleStreams {
			id, err := inventoryService.addArticleJob(ctx, articleData)
			if err != nil {
				logrus.WithError(err).Error("Something went wrong while addArticleJob for"+ articleData.Name)
				inventoryService.addArticleAsFailure(articleData, failedArticles, err)
				continue
			}
			inventoryService.addArticleAsSuccess(articleData, succedArticles, id)
		}
	waitChannel <- true
	}()
}

func (inventoryService *InventoryServiceImpl) addArticleJob(ctx context.Context, articleData *requests.Article) (int64, error) {
	article, err := inventoryService.getArticle(ctx,  articleData)
	if err != nil {
		if err.Error() == helpers.SQLRowNotFound {
			id, err := inventoryService.createNewArticle(ctx, articleData)
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
		updatedID, err := inventoryService.repo.UpdateArticle(ctx, article)
		if err != nil {
			return updatedID, err
		}
		return updatedID, nil
	}
	return -1, nil
}

func (inventoryService *InventoryServiceImpl) streamArticlesData(decoder *json.Decoder) chan *requests.Article {
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

func (inventoryService *InventoryServiceImpl) createNewArticle(ctx context.Context, value *requests.Article) (int64, error) {
	stock, err := strconv.Atoi(value.Stock)
	newArticle := database.NewArticle(0, stock, value.Name)
	id, err := inventoryService.repo.InsertArticle(ctx, newArticle)
	if err != nil {
		return id, err
	}
	return id, nil
}

func (inventoryService *InventoryServiceImpl)  getArticle(ctx context.Context,  value *requests.Article) (*database.Article, error) {
	return inventoryService.repo.GetArticle(ctx, helpers.GetArticleByName, value.Name)
}

func (inventoryService *InventoryServiceImpl) addArticleAsSuccess(article *requests.Article, suc []inventoryModel.SuccessfullyAddedArticle, id int64) {
	success := inventoryModel.SuccessfullyAddedArticle{}
	success.Name = article.Name
	success.Total = article.Stock
	success.Endpoint = "/api/v1/inventory/articles/" + strconv.Itoa(int(int64(id)))
	suc = append(suc, success)
}

func (inventoryService *InventoryServiceImpl) addArticleAsFailure(article *requests.Article, errorsArr []inventoryModel.FailedArticle, err error) {
	error := inventoryModel.FailedArticle{}
	error.Name = article.Name
	error.Reason = err.Error()
	errorsArr = append(errorsArr, error)
}