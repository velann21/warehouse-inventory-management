package requests

import (
	"encoding/json"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"io"
	"mime/multipart"
	"net/http"
)

const INVENTORYFILENAME = "inventory.json"

type AddArticles struct {
	Articles []Article `json:"articles"`
}

type Article struct {
	ArtID string `json:"art_id"`
	Name  string `json:"name"`
	Stock string `json:"stock"`
}

func NewAddArticles() *AddArticles {
	return &AddArticles{}
}

func (addArticles *AddArticles) ValidateAddArticles() error {
	if len(addArticles.Articles) < 0 {
		return helpers.InvalidRequest
	}
	return nil
}

func (addArticles *AddArticles) PopulateAddArticles(body io.Reader) error {
	decode := json.NewDecoder(body)
	err := decode.Decode(&addArticles)
	if err != nil {
		return helpers.InvalidRequest
	}
	return nil
}

func (addArticles *AddArticles) PopulateAddArticlesDataFromFile(req *http.Request, helper helpers.HelperBase) (*json.Decoder, *multipart.FileHeader, error) {
	err := req.ParseMultipartForm(20 << 20)
	if err != nil {
		return nil, nil, err
	}
	file, handler, err := req.FormFile("file")
	if err != nil {
		return nil, nil, err
	}
	defer file.Close()
	decode, err := helper.StreamFile(file)
	if err != nil {
		return nil, nil, err
	}
	return decode, handler, nil
}

func (addArticles *AddArticles) ValidateAddArticlesDataFromFile(req *http.Request, handler *multipart.FileHeader) error {
	if handler.Size > 20*1024*1024 {
		return helpers.InvalidRequest
	}
	if handler.Filename != INVENTORYFILENAME {
		return helpers.InvalidRequest
	}
	return nil
}
