package response

import (
	"encoding/json"
	"net/http"
)

type SuccessResponse struct {
	Success bool
	Status string
	Data []map[string]interface{}
}

func (entity *SuccessResponse) UserRegistrationResp(id *string) {
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["id"] = *id
	responseData = append(responseData, data)
	entity.Data = responseData
	entity.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "User registered"
}

func (entity *SuccessResponse) CreatePermissionResp(id *string){
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["id"] = *id
	responseData = append(responseData, data)
	entity.Data = responseData
	entity.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Permission Created"
}

func (entity *SuccessResponse) CreateRolesResp(id *string){
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["id"] = *id
	responseData = append(responseData, data)
	entity.Data = responseData
	entity.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Roles Created"
}

func (resp *SuccessResponse) CreateClusterResponse(id string, boo bool){
	responseData := make([]map[string]interface{}, 0)
	data := make(map[string]interface{})
	data["id"] = id
	data["Accepted"] = boo
	responseData = append(responseData, data)
	resp.Data = responseData
	resp.Success = true
	metaData := make(map[string]interface{})
	metaData["message"] = "Cluster queued to create in background, Please wait for 20+ min..."
}


func (resp *SuccessResponse) SuccessResponse(rw http.ResponseWriter, statusCode int){
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




