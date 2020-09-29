package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
	"github.com/velann21/warehouse-inventory-management/pkg/models/requests"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

type MockProductService struct {
	SrvType string 
}

func (productsService *MockProductService) AddProducts(ctx context.Context, products *requests.AddProducts) {
}

func (productsService *MockProductService) AddProductsFromFile(ctx context.Context, decoder *json.Decoder, waitChannel chan bool) error {
	if decoder == nil {
		return nil
	}
	return nil
}

func (productsService *MockProductService) PurchaseProducts(ctx context.Context, product *requests.PurchaseProduct) error {

	return nil
}

func (productsService *MockProductService) DeleteByID(ctx context.Context, productID int) (int64, error) {
	return -1, nil
}

func (productsService *MockProductService) GetProductByID(ctx context.Context, productDetails *requests.GetProductDetails) ([]*database.ProductDetails, error) {
	if productsService.SrvType == ""{
		fmt.Println("Mock details")
		details := []*database.ProductDetails{}
		return details, nil
	}else if productsService.SrvType == "200Response" {
		details := []*database.ProductDetails{
			&database.ProductDetails{
			QuantityEach:1,
			ArtID:10,
			Name:"legs",
			},
			&database.ProductDetails{
				QuantityEach:20,
				ArtID:11,
				Name:"Bolts",
			},
		}
		return details, nil
	}
	return nil, nil
}

func (productsService *MockProductService) GetAllProducts(ctx context.Context) ([]*database.Product, error) {

	return nil, nil
}

func Router(prodService services.ProductsService) *mux.Router {
	router := mux.NewRouter()
	ctrlObj := ProductsControllers{service: prodService}
	router.HandleFunc("/v1/inventory/products", ctrlObj.AddProducts).Methods("POST")
	router.HandleFunc("/v1/inventory/products/{id}", ctrlObj.GetProductDetails).Methods("POST")
	return router
}

func TestProductsControllers_AddProducts(t *testing.T) {
	req := `{
  "products": [
    {
      "name": "Dining Table",
      "contain_articles": [
        {
          "art_id": "1",
          "amount_of": "4"
        },
        {
          "art_id": "1",
          "amount_of": "8"
        },
        {
          "art_id": "3",
          "amount_of": "1"
        }
        ]
    }
  ]
}
`

	request, _ := http.NewRequest(http.MethodPost, "/v1/inventory/products", strings.NewReader(req))
	response := httptest.NewRecorder()
	Router(&MockProductService{}).ServeHTTP(response, request)
	fmt.Println(response.Body)
}

func TestProductsControllers_GetProductDetailsWithEmptyResponse(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/v1/inventory/products/15", nil)
	response := httptest.NewRecorder()
	Router(&MockProductService{}).ServeHTTP(response, request)
	fmt.Println(response.Body)
}

func TestProductsControllers_GetProductDetailsWith200Response(t *testing.T) {
	request, _ := http.NewRequest(http.MethodPost, "/v1/inventory/products/15", nil)
	response := httptest.NewRecorder()
	Router(&MockProductService{SrvType:"200Response"}).ServeHTTP(response, request)
	fmt.Println(response.Body)
}
