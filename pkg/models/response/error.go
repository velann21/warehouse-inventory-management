package response

import (
	"encoding/json"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"net/http"
)

type ErrorResponse struct {
	Success bool
	Errors  []Error
}

func NewErrorResponse()*ErrorResponse{
	return &ErrorResponse{}
}

type Error struct {
	Message   string
	ErrorCode int
}

func (err *ErrorResponse) HandleError(er error, w http.ResponseWriter) {
	if er == nil {
		logrus.Error("invalid error")
		return
	}
	errList := make([]Error, 0)
	switch er {
	case helpers.InvalidRequest:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 1,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(400)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)

	case helpers.SomethingWrong:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 1,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	default:
		errObj := Error{
			Message:   er.Error(),
			ErrorCode: 1,
		}
		errList = append(errList, errObj)
		resp := ErrorResponse{
			Success: false,
			Errors:  errList,
		}
		w.WriteHeader(500)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(resp)
	}

}
