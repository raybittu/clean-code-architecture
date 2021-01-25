package delivery

import "net/http"

type Customer interface {
	GetById(w http.ResponseWriter, r *http.Request)
	GetByName(w http.ResponseWriter, r *http.Request)
	PostCustomer(w http.ResponseWriter, r *http.Request)
	PutCustomer(w http.ResponseWriter, r *http.Request)
	DeleteCustomer(w http.ResponseWriter, r *http.Request)
}
