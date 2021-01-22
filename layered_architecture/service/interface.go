package service

import (
	"layered/architecture/entities"
	"net/http"
)

type Customer interface {
	GetByID(w http.ResponseWriter, id int)
	GetByName(w http.ResponseWriter, name string)
	CreateCustomer(w http.ResponseWriter, c entities.Customer)
	UpdateCustomer(w http.ResponseWriter, id int, c entities.Customer)
	DeleteCustomer(w http.ResponseWriter, id int)
}
