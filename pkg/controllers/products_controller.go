package controllers

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/models/response"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
	"net/http"
)

type ProductsControllers struct {
	service services.ProductsService
	helper  helpers.HelperBase
}

func NewProductsController(service services.ProductsService, helpers helpers.HelperBase) *ProductsControllers {
	return &ProductsControllers{service: service, helper: helpers}
}

func (productsControllers *ProductsControllers) AddProducts(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the AddProducts().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	products := requests.NewAddProducts()
	err := products.PopulateAddProducts(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	err = products.ValidateAddProducts()
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	productsControllers.service.AddProducts(ctx, products)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done AddProducts().")
	return
}

func (productsControllers *ProductsControllers) AddProductsFromFile(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the AddProductsFromFile().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	products := requests.NewAddProducts()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	decode, handler, err := products.PopulateAddProductsDataFromFile(req, productsControllers.helper)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateAddProductsDataFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}
	err = products.ValidateAddProductsDataFromFile(req, handler)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateAddArticles() ")
		errorResponse.HandleError(err, res)
		return
	}
	waitChannel := make(chan bool)
	err = productsControllers.service.AddProductsFromFile(ctx, decode, waitChannel)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while AddProductsFromFile() ")
		errorResponse.HandleError(err, res)
		return
	}

	<-waitChannel
	close(waitChannel)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done AddProductsFromFile().")
	return
}

func (productsControllers *ProductsControllers) ListProducts(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the GetAllProducts().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()

	productList, err := productsControllers.service.GetAllProducts(ctx)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while GetAllProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	successResponse.GetProductsResponse(productList)
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done ListProducts().")
	return
}

func (productsControllers *ProductsControllers) PurchaseProducts(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the BuyProducts().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	purchseProducts := requests.NewPurchaseProducts()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	err := purchseProducts.PopulatePurchaseProducts(req.Body)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulatePurchaseProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	err = purchseProducts.ValidatePurchaseProducts()
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidatePurchaseProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	err = productsControllers.service.PurchaseProducts(ctx, purchseProducts)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PurchaseProducts()")
		errorResponse.HandleError(err, res)
		return
	}
	successResponse.SuccessResponse(res, http.StatusOK)
	logrus.Info("Done PurchaseProducts().")
	return
}

func (productsControllers *ProductsControllers) GetProductDetails(res http.ResponseWriter, req *http.Request) {
	logrus.Info("Starting the GetProductDetails().......")
	ctx, cancel := context.WithTimeout(req.Context(), TIMEOUT)
	defer cancel()
	getProductDetails := requests.NewGetProductDetails()
	successResponse := response.NewSuccessResponse()
	errorResponse := response.NewErrorResponse()
	id := mux.Vars(req)["id"]
	err := getProductDetails.PopulateGetProductDetails(id)
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while PopulateGetProductDetails()")
		errorResponse.HandleError(err, res)
		return
	}
	err = getProductDetails.ValidateGetProductDetails()
	if err != nil {
		logrus.WithError(err).Error("Something went wrong while ValidateGetProductDetails()")
		errorResponse.HandleError(err, res)
		return
	}

	productDetail, err := productsControllers.service.GetProductByID(ctx, getProductDetails)
	if err != nil{
		logrus.WithError(err).Error("Something went wrong while GetProductByID()")
		errorResponse.HandleError(err, res)
		return
	}
	successResponse.GetProductDetailsResponse(productDetail)
	successResponse.SuccessResponse(res, http.StatusOK)
	return
}
