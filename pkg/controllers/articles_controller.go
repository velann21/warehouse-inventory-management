package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/models/response"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
	"net/http"
	"time"
)

const TIMEOUT = time.Second * 10

//1. Todo make the response structure into Hateos style
//2. Todo change the database table name from inventory to articles
type ArticlesControllers struct {
	service services.ArticlesService
	helper  helpers.HelperBase
}

func NewArticlesController(service services.ArticlesService, helpers helpers.HelperBase) *ArticlesControllers {
	return &ArticlesControllers{service: service, helper: helpers}
}

func (productsControllers *ArticlesControllers) AddArticles(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the AddArticles().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	articles := requests.NewAddArticles()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	err := articles.PopulateAddArticles(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateAddArticles() ")
		errorResponse.HandleError(err, res)
		return
	}
	err = articles.ValidateAddArticles()
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddArticles() ")
		errorResponse.HandleError(err, res)
		return
	}
	addedArticles, failedArticles := productsControllers.service.AddArticles(ctx, articles)
	successResponse.AddArticlesResponse(addedArticles, failedArticles)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done AddArticles().")
	return
}

func (productsControllers *ArticlesControllers) AddArticlesFromFile(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the AddArticlesFromFile().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	articles := requests.NewAddArticles()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	decode, handler, err := articles.PopulateAddArticlesDataFromFile(req, productsControllers.helper)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateAddArticlesDataFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}
	err = articles.ValidateAddArticlesDataFromFile(req, handler)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddArticlesDataFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}
	waitChannel := make(chan bool)
	err = productsControllers.service.AddArticlesFromFile(ctx, decode, waitChannel)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while AddArticlesFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}
	<-waitChannel
	close(waitChannel)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done AddArticlesFromFile().")
	return
}

func (productsControllers *ArticlesControllers) ListArticles(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the ListArticles().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	articleList, err := productsControllers.service.ListArticles(ctx)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ListArticles() ")
		errorResponse.HandleError(err, res)
		return
	}
	successResponse.GetArticlesResponse(articleList)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done ListArticles().")
	return
}

func (productsControllers *ArticlesControllers) SqlMigration(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the SqlMigration().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()

	sqlMigrationReq := requests.SqlMigrationRequest{}
	err := sqlMigrationReq.PopulateSqlMigrationRequest(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateSqlMigrationRequest() ")
		errorResponse.HandleError(err, res)
		return
	}
	err = sqlMigrationReq.ValidateSqlMigrationRequest()
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateSqlMigrationRequest() ")
		errorResponse.HandleError(err, res)
		return
	}

	err = productsControllers.service.SqlMigration(ctx, &sqlMigrationReq)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while SqlMigration() ")
		errorResponse.HandleError(err, res)
		return
	}
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done the SqlMigration().......")
	return
}
