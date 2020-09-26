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

const TIMEOUT = time.Second * 50

type InventoryControllers struct {
	service services.InventoryService
	helper  helpers.HelperBase
}

func NewInventoryController(service services.InventoryService, helpers helpers.HelperBase) *InventoryControllers {
	return &InventoryControllers{service: service, helper: helpers}
}

func (inventoryControllers *InventoryControllers) AddArticles(res http.ResponseWriter, req *http.Request) {
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
	addedArticles, failedArticles := inventoryControllers.service.AddArticles(ctx, articles)
	successResponse.AddArticlesResponse(addedArticles, failedArticles)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done AddArticles().")
	return
}

func (inventoryControllers *InventoryControllers) AddArticlesFromFile(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the AddArticlesFromFile().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	articles := requests.NewAddArticles()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	decode,handler, err := articles.PopulateAddArticlesDataFromFile(req, inventoryControllers.helper)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateAddArticlesDataFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}
	err = articles.ValidateAddArticlesDataFromFile(req, handler)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddArticles() ")
		errorResponse.HandleError(err, res)
		return
	}
	waitChannel := make(chan bool)
	err = inventoryControllers.service.AddArticlesFromFile(ctx, decode, waitChannel)
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

func (inventoryControllers *InventoryControllers) GetArticles(res http.ResponseWriter, req *http.Request) {

}

func (inventoryControllers *InventoryControllers) AddProducts(res http.ResponseWriter, req *http.Request) {
	logrus.Info("asasas")
}

func (inventoryControllers *InventoryControllers) GetProducts(res http.ResponseWriter, req *http.Request) {

}