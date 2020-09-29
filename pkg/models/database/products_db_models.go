package database

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name, omitempty"`
	Price       string `json:"price, omitempty"`
	Description string `json:"description, omitempty"`
	Quantity    int    `json:"available_quantity, omitempty"`
}

func NewProducts(ID int, name, description, price string, quantity int) *Product {
	return &Product{ID: ID, Name: name, Price: price, Description: description, Quantity: quantity}
}

func (product *Product) SetName(name string) {
	product.Name = name
}

func (product *Product) SetDescription(description string) {
	product.Description = description
}

func (product *Product) SetQuantity(quantity int) {
	product.Quantity = quantity
}

func (product *Product) GetPrice() string {
	return product.Price
}

func (product *Product) GetName() string {
	return product.Name
}

func (product *Product) GetDescription() string {
	return product.Description
}

func (product *Product) GetQuantity() int {
	return product.Quantity
}

type ProductsArticles struct {
	ProductID             int `json:"product_id"`
	InventoryID           int `json:"inventory_id"`
	EachQuantity          int `json:"quantity_each"`
	TotalRequiredQuantity int `json:"total_required_quantity"`
}

func NewProductsArticles(productID int, inventoryID int, eachQuantity int, totalRequiredQuantity int) *ProductsArticles {
	return &ProductsArticles{ProductID: productID, InventoryID: inventoryID, EachQuantity: eachQuantity, TotalRequiredQuantity: totalRequiredQuantity}
}

func (productArticles *ProductsArticles) GetProductID() int {
	return productArticles.ProductID
}

func (productArticles *ProductsArticles) GetInventoryID() int {
	return productArticles.InventoryID
}

func (productArticles *ProductsArticles) GetEachQuantity() int {
	return productArticles.EachQuantity
}

func (productArticles *ProductsArticles) GetTotalRequiredQuantity() int {
	return productArticles.TotalRequiredQuantity
}

type ProductDetails struct {
	QuantityEach int    `json:"quantity_each"`
	ArtID        int    `json:"art_id"`
	Name         string `json:"name"`
}
