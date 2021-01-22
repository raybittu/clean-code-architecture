package service

import (
	"layered/architecture/entities"
	"net/http"
)

type Customer interface {
	GetByID(w http.ResponseWriter, id int)
	GetByName(w http.ResponseWriter, name string)
	CreateCustomer(w http.ResponseWriter, cust entities.Customer)
	UpadteCustomer(w http.ResponseWriter, id int, customer entities.Customer)
	DeleteCustomer(w http.ResponseWriter, id int)
}
