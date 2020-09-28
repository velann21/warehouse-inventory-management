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
	helper := helpers.NewHelper(helpers.HELPER_VERSION_V1)
	helper.SetEnv()

	sqlconn := db.NewSqlConnection()
	sqlConn, err := sqlconn.OpenSqlConnection(helper)
	if err != nil {
		logrus.WithField("EventType", "DbConnection").WithError(err).Error("Db Connection Error")
		os.Exit(100)
	}

	err = sqlConn.PingContext(context.Background())
	if err != nil {
		logrus.WithField("EventType", "PingContext").WithError(err).Error("Mysql PingContext Error")
		os.Exit(100)
	}

	r := mux.NewRouter().StrictSlash(false)
	mainRoutes := r.PathPrefix("/api").Subrouter()

	articleRoutes := routes.NewArticlesRoutes(sqlConn, helper)
	articleRoutes.ArticlesRoutes(mainRoutes)

	productRoutes := routes.NewProductsRoutes(sqlConn, helper)
	productRoutes.ProductRoutes(mainRoutes)

	logrus.Info("Starting the server with port :8083")
	if err := http.ListenAndServe(":8083", r); err != nil {
		logrus.WithField("EventType", "Server Bootup").WithError(err).Error("Server Bootup Error")
		log.Fatal(err)
		return
	}
}
