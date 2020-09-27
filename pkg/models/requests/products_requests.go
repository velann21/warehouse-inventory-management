package requests

import (
	"encoding/json"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"io"
	"mime/multipart"
	"net/http"
)

const PRODUCTSFILENAME = "products.json"

type AddProducts struct {
	Products []Products `json:"products"`
}

type Products struct {
	Name     string            `json:"name"`
	Price    string            `json:"price"`
	Articles []ProductArticles `json:"contain_articles"`
}

type ProductArticles struct {
	ArtID    string `json:"art_id"`
	AmountOf string `json:"amount_of"`
}

func NewAddProducts() *AddProducts {
	return &AddProducts{}
}

func (products *AddProducts) PopulateAddProducts(body io.Reader) error {
	decode := json.NewDecoder(body)
	err := decode.Decode(&products)
	if err != nil {
		return helpers.InvalidRequest
	}
	return nil
}

func (products *AddProducts) ValidateAddProducts() error {
	if len(products.Products) <= 0 {
		return helpers.InvalidRequest
	}
	return nil
}

func (addProducts *AddProducts) PopulateAddProductsDataFromFile(req *http.Request, helper helpers.HelperBase) (*json.Decoder, *multipart.FileHeader, error) {
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

func (addProducts *AddProducts) ValidateAddProductsDataFromFile(req *http.Request, handler *multipart.FileHeader) error {
	if handler.Size > 20*1024*1024 {
		return helpers.InvalidRequest
	}
	if handler.Filename != PRODUCTSFILENAME {
		return helpers.InvalidRequest
	}
	return nil
}

type PurchaseProduct struct {
	Name string `json:"name"`
	ID   string `json:"id"`
}

func NewPurchaseProducts() *PurchaseProduct {
	return &PurchaseProduct{}
}
func (purchaseProduct *PurchaseProduct) PopulatePurchaseProducts(body io.Reader) error {
	decode := json.NewDecoder(body)
	err := decode.Decode(&purchaseProduct)
	if err != nil {
		return helpers.InvalidRequest
	}
	return nil
}

func (purchaseProduct *PurchaseProduct) ValidatePurchaseProducts() error {

	if purchaseProduct.Name == "" {
		return helpers.InvalidRequest
	}
	if purchaseProduct.ID == "" {
		return helpers.InvalidRequest
	}

	return nil
}

type GetProductDetails struct {
	ID string `json:"id"`
}

func NewGetProductDetails() *GetProductDetails {
	return &GetProductDetails{}
}
func (getProductDetails *GetProductDetails) PopulateGetProductDetails(ID string) error {
	getProductDetails.ID = ID
	return nil
}

func (getProductDetails *GetProductDetails) ValidateGetProductDetails() error {
	if getProductDetails.ID == "" {
		return helpers.InvalidRequest
	}

	return nil
}
