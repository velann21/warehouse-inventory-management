package main

import (
	"context"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
	db "github.com/velann21/warehouse-inventory-management/pkg/databases"
	"github.com/velann21/warehouse-inventory-management/pkg/helpers"
	"github.com/velann21/warehouse-inventory-management/pkg/routes"
	"log"
	"net/http"
	"os"
)

func main() {

	helpers.SetEnv()
	sqlconn := db.SQLConnection{}

	err := sqlconn.OpenSqlConnection()
	if err != nil {
		logrus.WithField("EventType", "DbConnection").WithError(err).Error("Db Connection Error")
		os.Exit(100)
	}

	err = db.GetSqlConnection().PingContext(context.Background())
	if err != nil {
		fmt.Println(err)
		os.Exit(100)
	}

	r := mux.NewRouter().StrictSlash(false)
	rmRoutes := routes.InventoryRoutes{}
	rmRoutes.InventoryRoutes(r)

	if err := http.ListenAndServe(":8080", r); err != nil {
		logrus.WithField("EventType", "Server Bootup").WithError(err).Error("Server Bootup Error")
		log.Fatal(err)
		return
	}
}
