package main

import (
	"context"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	db "github.com/velann21/warehouse-inventory-management/pkg/helpers/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/routes"
	"log"
	"net/http"
	"os"
)

func main() {
	//Helpers Object
	helper := helpers.NewHelper(helpers.HELPER_VERSION_V1)
	helper.SetEnv()

	//Sql Connection Object, If failed to make sql conn restart container/app
	sqlconn := db.NewSqlConnection()
	sqlConn, err := sqlconn.OpenSqlConnection(helper)
	if err != nil {
		logrus.WithField("EventType", "DbConnection").WithError(err).Error("Db Connection Error")
		os.Exit(100)
	}

	// To make sure sql is up if not restart container/app
	err = sqlConn.PingContext(context.Background())
	if err != nil {
		logrus.WithField("EventType", "PingContext").WithError(err).Error("Mysql PingContext Error")
		os.Exit(100)
	}

	r := mux.NewRouter().StrictSlash(false)
	mainRoutes := r.PathPrefix("/api").Subrouter()

	// ArticlesRoutes Object Init
	articleRoutes := routes.NewArticlesRoutes(sqlConn, helper)
	articleRoutes.ArticlesRoutes(mainRoutes)

	// ProductsRoutes Object Init
	productRoutes := routes.NewProductsRoutes(sqlConn, helper)
	productRoutes.ProductRoutes(mainRoutes)

	//Starting server and Listen on port 8083
	logrus.Info("Starting the server with port :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		logrus.WithField("EventType", "Server Bootup").WithError(err).Error("Server Bootup Error")
		log.Fatal(err)
		return
	}
}
