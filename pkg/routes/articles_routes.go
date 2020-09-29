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

type ArticlesRoutes struct {
	controller *controllers.ArticlesControllers
}

func NewArticlesRoutes(sql *sql.DB, helpers helpers.HelperBase) *ArticlesRoutes {
	// Initializing all the dependent object for Articles at one place
	sqlClient := databases.NewSqlClient(sql)
	inventoryRepository := repository.NewArticlesRepositoryFactory(repository.ARTICLES_REPO_VERSION1, sqlClient)
	inventoryService := services.NewInventoryServiceFactory(services.ARTICLES_SERVICE_VERSION1, inventoryRepository)
	return &ArticlesRoutes{controller: controllers.NewArticlesController(inventoryService, helpers)}
}

func (inventory *ArticlesRoutes) ArticlesRoutes(route *mux.Router) {
	route.Path("/v1/inventory/articles").HandlerFunc(inventory.controller.AddArticles).Methods("POST")
	route.Path("/v1/inventory/articles/fromFile").HandlerFunc(inventory.controller.AddArticlesFromFile).Methods("POST")
	route.Path("/v1/inventory/articles").HandlerFunc(inventory.controller.ListArticles).Methods("GET")

	// Todo Move this to common routes folders later
	route.PathPrefix("/v1/inventory/sqlmigration").HandlerFunc(inventory.controller.SqlMigration).Methods("POST")
}
