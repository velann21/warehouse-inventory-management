package requests

import (
	"encoding/json"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"io"
)

type AddArticles struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	Name string `json:"name"`
	Stock int `json:"stock"`
}

func (addArticles *AddArticles) ValidateAddArticles()error{
	if len(addArticles.Articles) < 0{
		return helpers.InvalidRequest
	}
	return nil
}

func (addArticles *AddArticles) PopulateAddArticles(body io.Reader)error{
	decode := json.NewDecoder(body)
	err := decode.Decode(&addArticles)
	if err != nil{
		return helpers.InvalidRequest
	}
	return nil
}
