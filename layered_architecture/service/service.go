package service

import "C"
import (
	"encoding/json"
	"layered/architecture/entities"
	"layered/architecture/store"
	"net/http"
	"strings"
	"time"
)

type CustomerService struct {
	store store.Customer
}

func New(customer store.Customer) CustomerService {
	return CustomerService{store: customer}
}

func (c CustomerService) GetById(w http.ResponseWriter, id int) {
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("id can't be less than 1"))
		return
	}
	resp, err := c.store.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode([]entities.Customer(nil))
	} else {
		if resp.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			json.NewEncoder(w).Encode([]entities.Customer(nil))
		} else {
			json.NewEncoder(w).Encode(resp)
		}
	}
}

func (c CustomerService) GetByName(w http.ResponseWriter, name string) {
	if len(name) <= 0 {
		resp, err := c.store.GetByName("")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(resp)
			return
		}
		json.NewEncoder(w).Encode(resp)
		return
	}
	resp, err := c.store.GetByName(name)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(resp)
		return
	}
	if len(resp) == 0 {
		w.WriteHeader(http.StatusNotFound)
	}
	json.NewEncoder(w).Encode(resp)
}

func (c CustomerService) CreateCustomer(w http.ResponseWriter, cust entities.Customer) {

	if cust.Name == "" || cust.DOB == "" || cust.Address.StreetName == "" || cust.Address.City == "" || cust.Address.State == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("some fields are missing"))
	} else if timestamp := DateSubstract(cust.DOB); timestamp/(3600*24*12*30) < 18 {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("You are below 18, so you are not allowed to be our customer"))
	} else {
		cust, err := c.store.Create(cust)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("something went wrong"))
			return
		}
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(cust)
	}
}

func DateSubstract(d1 string) int {
	d1_slice := strings.Split(d1, "/")

	newDate := d1_slice[2] + "-" + d1_slice[1] + "-" + d1_slice[0]
	myDate, err := time.Parse("2006-01-02", newDate)

	if err != nil {
		//panic(err)
		return 0
	}

	return int(time.Now().Unix() - myDate.Unix())
}
