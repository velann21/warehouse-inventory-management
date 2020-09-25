package routes

import (
	"database/sql"
	"github.com/gorilla/mux"
	"github.com/velann21/warehouse-inventory-management/pkg/controllers"
	"github.com/velann21/warehouse-inventory-management/pkg/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/repository"
	"github.com/velann21/warehouse-inventory-management/pkg/services"
)

type InventoryRoutes struct {
	controller *controllers.InventoryControllers
}

func NewRoutes(sql *sql.DB)*InventoryRoutes{
	sqlClient := databases.NewSqlClient(sql)
	inventoryRepository := repository.NewInventoryRepositoryFactory(repository.INVENTORY_REPO_VERSION1, sqlClient)
	inventoryService := services.NewInventoryServiceFactory(services.SERVICE_VERSION1, inventoryRepository)
	return &InventoryRoutes{controller:controllers.NewInventoryController(inventoryService)}
}

func (inventory *InventoryRoutes) InventoryRoutes(route *mux.Router) {
	route.PathPrefix("/v1/inventory/articles").HandlerFunc(inventory.controller.AddArticles).Methods("POST")
	route.PathPrefix("/v1/inventory/articles").HandlerFunc(inventory.controller.GetArticles).Methods("GET")
	route.PathPrefix("/v1/inventory/products").HandlerFunc(inventory.controller.AddProducts).Methods("POST")
	route.PathPrefix("/v1/inventory/products").HandlerFunc(inventory.controller.GetProducts).Methods("GET")
}
