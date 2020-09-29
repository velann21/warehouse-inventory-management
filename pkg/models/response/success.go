package response

import (
	"encoding/json"
	"github.com/velann21/warehouse-inventory-management/pkg/models/database"
	inventoryArticle "github.com/velann21/warehouse-inventory-management/pkg/models/internals"
	"net/http"
)

type SuccessResponse struct {
	Success bool
	Status  string
	Data    []map[string]interface{}
}

func NewSuccessResponse() *SuccessResponse {
	return &SuccessResponse{}
}

func (resp *SuccessResponse) AddArticlesResponse(successArticle []inventoryArticle.SuccessfullyAddedArticle, failedArticle []inventoryArticle.FailedArticle) {
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["addedArticles"] = successArticle
	data["failedArticles"] = failedArticle
	responseData = append(responseData, data)
	resp.Data = responseData
	resp.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Add article job is completed"
}

func (resp *SuccessResponse) GetProductsResponse(products []*database.Product) {
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["products"] = products
	responseData = append(responseData, data)
	resp.Data = responseData
	resp.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Success"
}

func (resp *SuccessResponse) GetArticlesResponse(articles []*database.Article) {
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["articles"] = articles
	responseData = append(responseData, data)
	resp.Data = responseData
	resp.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Success"
}

func (resp *SuccessResponse) GetProductDetailsResponse(productDetails []*database.ProductDetails) {
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["product_details"] = productDetails
	responseData = append(responseData, data)
	resp.Data = responseData
	resp.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Success"
}

func (resp *SuccessResponse) SuccessResponse(rw http.ResponseWriter, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")

	switch statusCode {
	case http.StatusOK:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
		resp.Success = true
	case http.StatusCreated:
		rw.WriteHeader(http.StatusCreated)
		resp.Status = http.StatusText(http.StatusCreated)
		resp.Success = true
	case http.StatusAccepted:
		rw.WriteHeader(http.StatusAccepted)
		resp.Status = http.StatusText(http.StatusAccepted)
		resp.Success = true
	default:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
		resp.Success = true
	}
	// send response
	_ = json.NewEncoder(rw).Encode(resp)
	return
}
