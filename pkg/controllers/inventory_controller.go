package controllers

import (
	"context"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
	"net/http"
	"time"
)

const TIMEOUT = time.Second*10
type InventoryControllers struct {
	service services.Service
}

func NewInventoryController(service services.Service)*InventoryControllers{
	return &InventoryControllers{service:service}
}

func (inventoryControllers *InventoryControllers) AddArticles(res http.ResponseWriter, req *http.Request){
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	articles := requests.AddArticles{}
	err := articles.PopulateAddArticles(req.Body)
	if err != nil{
		logrus.WithError(err).Error("Something wrong while PopulateAddArticles: ")
		return
	}
	err = articles.ValidateAddArticles()
	if err != nil{
		logrus.WithError(err).Error("Something went wrong while ValidateAddArticles: ")
		return
	}

	err = inventoryControllers.service.AddArticles(ctx, articles)
	if err != nil{
		logrus.WithError(err).Error("Something went wrong while AddArticles: ")
		return
	}
	return
}

func (inventoryControllers *InventoryControllers) GetArticles(res http.ResponseWriter, req *http.Request){

}

func (inventoryControllers *InventoryControllers) AddProducts(res http.ResponseWriter, req *http.Request){

}

func (inventoryControllers *InventoryControllers) GetProducts(res http.ResponseWriter, req *http.Request){

}
