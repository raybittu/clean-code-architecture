package main

import (
	"database/sql"
	"log"
	"net/http"

	"layered/architecture/delivery"
	"layered/architecture/service"
	"layered/architecture/store"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	var db, dbErr = sql.Open("mysql", "root:1118209@/Customer_service")
	if dbErr != nil {
		panic(dbErr)
	}
	defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr.Error()) // proper error handling instead of panic in your app
	}
	r := mux.NewRouter()
	datastore := store.New(db)
	service := service.New(datastore)
	handler := delivery.New(service)

	r.HandleFunc("/customer", handler.GetByName).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", handler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/customer", handler.PostCustomer).Methods(http.MethodPost)
	log.Fatal(http.ListenAndServe(":8080", r))
}
