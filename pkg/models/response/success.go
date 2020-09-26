package response

import (
	"encoding/json"
	inventoryArticle "github.com/velann21/warehouse-inventory-management/pkg/models/internals"
	"net/http"
)

type SuccessResponse struct {
	Success bool
	Status  string
	Data    []map[string]interface{}
}

func NewSuccessResponse()*SuccessResponse{
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
func (resp *SuccessResponse) SuccessResponse(rw http.ResponseWriter, statusCode int) {
	rw.Header().Set("Content-Type", "application/json")

	switch statusCode {
	case http.StatusOK:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
	case http.StatusCreated:
		rw.WriteHeader(http.StatusCreated)
		resp.Status = http.StatusText(http.StatusCreated)
	case http.StatusAccepted:
		rw.WriteHeader(http.StatusAccepted)
		resp.Status = http.StatusText(http.StatusAccepted)
	default:
		rw.WriteHeader(http.StatusOK)
		resp.Status = http.StatusText(http.StatusOK)
	}
	// send response
	_ = json.NewEncoder(rw).Encode(resp)
	return
}
