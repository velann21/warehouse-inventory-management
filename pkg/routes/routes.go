package routes

import "github.com/gorilla/mux"

type InventoryRoutes struct {
}

func (i *InventoryRoutes) InventoryRoutes(route *mux.Router) {
	route.PathPrefix("/inventory/articles").HandlerFunc().Methods("POST")
	route.PathPrefix("/inventory/articles").HandlerFunc().Methods("GET")
	route.PathPrefix("/inventory/products").HandlerFunc().Methods("POST")
	route.PathPrefix("/inventory/products").HandlerFunc().Methods("GET")
}
