package main

import (
	"log"
	"net/http"

	"layered/architecture/delivery"
	"layered/architecture/service"
	"layered/architecture/store"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

func main() {
	r := mux.NewRouter()
	datastore := store.New()
	defer datastore.Close()
	service := service.New(datastore)
	handler := delivery.New(service)

	r.HandleFunc("/customer", handler.GetByName).Methods(http.MethodGet)
	r.HandleFunc("/customer/{id}", handler.GetById).Methods(http.MethodGet)
	r.HandleFunc("/customer", handler.PostCustomer).Methods(http.MethodPost)
	r.HandleFunc("/customer/{id}", handler.PutCustomer).Methods(http.MethodPut)
	r.HandleFunc("/customer/{id}", handler.DeleteCustomer).Methods(http.MethodDelete)
	log.Fatal(http.ListenAndServe(":8080", r))
}
