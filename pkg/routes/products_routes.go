package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/velann21/warehouse-inventory-management/pkg/controllers"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
)

type ProductsRoutes struct {
	controller *controllers.ProductsControllers
}

func NewProductsRoutes(sql *sql.DB, helpers helpers.HelperBase) *ProductsRoutes {
	sqlClient := databases.NewSqlClient(sql)
	inventoryRepository := repository.NewProductsRepositoryFactory(repository.PRODUCTS_REPO_VERSION1, sqlClient)
	productsService := services.NewProductsServiceFactory(services.PRODUCTS_SERVICE_VERSION1, inventoryRepository)
	return &ProductsRoutes{controller: controllers.NewProductsController(productsService, helpers)}
}

func (productRoutes *ProductsRoutes) ProductRoutes(route *mux.Router) {
	route.PathPrefix("/v1/inventory/products/fromFile").HandlerFunc(productRoutes.controller.AddProductsFromFile).Methods("POST")
	route.PathPrefix("/v1/inventory/products").HandlerFunc(productRoutes.controller.AddProducts).Methods("POST")
	route.PathPrefix("/v1/inventory/products/{id}").HandlerFunc(productRoutes.controller.GetProductDetails).Methods("GET")
	route.PathPrefix("/v1/inventory/products").HandlerFunc(productRoutes.controller.ListProducts).Methods("GET")
	route.PathPrefix("/v1/inventory/purchaseProducts").HandlerFunc(productRoutes.controller.PurchaseProducts).Methods("POST")
}
