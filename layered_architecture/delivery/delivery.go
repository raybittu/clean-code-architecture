package delivery

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"layered/architecture/entities"
	"layered/architecture/errors"
	"layered/architecture/service"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type CustomerHandler struct {
	service service.Customer
}

func New(customer service.Customer) Customer {
	return CustomerHandler{service: customer}
}

func (c CustomerHandler) GetById(w http.ResponseWriter, r *http.Request) {
	pathparams := mux.Vars(r)

	id, err := strconv.Atoi(pathparams["id"])
	if err != nil {
		_, _ = w.Write([]byte(errors.ErrInvalidID))
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if id <= 0 {
		w.WriteHeader(http.StatusBadRequest)
		//err := customerError.MissingParams{Params: []string{"id"}}
		////fmt.Println(err)
		//data, err1 := json.Marshal(err.Error())
		//if err1 != nil {
		//	w.WriteHeader(http.StatusBadRequest)
		//}
		//_, _ = w.Write(data)
		w.Write([]byte(errors.ErrInvalidID))
		return
	}
	resp, err1 := c.service.GetByID(id)

	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err1.Error()))
		//_ = json.NewEncoder(w).Encode([]entities.Customer(nil))
	} else {
		if resp.ID == 0 {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(errors.ErrDBNotFound))
			//_ = json.NewEncoder(w).Encode([]entities.Customer(nil))
		} else {
			body, _ := json.Marshal(resp)
			_, _ = w.Write(body)
		}
	}

}

func (c CustomerHandler) GetByName(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	name, ok := params["name"]

	if !ok {
		resp, err1 := c.service.GetAll()
		if err1 != nil {
			w.WriteHeader(http.StatusInternalServerError)
			//_ = json.NewEncoder(w).Encode(resp)
			w.Write([]byte(err1.Error()))
			return
		}
		if len(resp) == 0 {
			w.WriteHeader(http.StatusNotFound)
			//_ = json.NewEncoder(w).Encode(resp)
			w.Write([]byte(errors.ErrDBNotFound))
			return
		}
	}
	resp, err1 := c.service.GetByName(name[0])

	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		//_ = json.NewEncoder(w).Encode(resp)
		w.Write([]byte(err1.Error()))
		return
	}
	if len(resp) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err1.Error()))
		return
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func (c CustomerHandler) PostCustomer(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	var cust entities.Customer
	err = json.Unmarshal(body, &cust)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid data format"))
		return
	}

	if cust.Name == "" || cust.DOB == "" || cust.Address.StreetName == "" || cust.Address.City == "" || cust.Address.State == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("some fields are missing"))
		return
	}
	cust, err1 := c.service.CreateCustomer(cust)
	fmt.Println(err1)
	if err1 != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(err1.Error()))
		return
	}
	w.WriteHeader(http.StatusCreated)
	_ = json.NewEncoder(w).Encode(cust)
}

func (c CustomerHandler) PutCustomer(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)["id"]
	id, err := strconv.Atoi(pathParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.ErrInvalidID))
		return
	}
	var customer entities.Customer
	bodyData, _ := ioutil.ReadAll(r.Body)
	err = json.Unmarshal(bodyData, &customer)

	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte(errors.ErrJSONFormat))
		return
	}
	if customer.ID != 0 || customer.DOB != "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("Can't update Id or DOB"))
		return
	}
	if customer.Name == "" && customer.Address.State == "" && customer.Address.City == "" && customer.Address.StreetName == "" {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("No data to update"))
		return
	}
	cust, err1 := c.service.UpdateCustomer(id, customer)

	if err1 != nil {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(err1.Error()))
		return
	}
	_ = json.NewEncoder(w).Encode(cust)

}

func (c CustomerHandler) DeleteCustomer(w http.ResponseWriter, r *http.Request) {
	pathParams := mux.Vars(r)["id"]

	id, err := strconv.Atoi(pathParams)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		_, _ = w.Write([]byte("invalid id parameter"))
		return
	}
	resp, err1 := c.service.DeleteCustomer(id)
	if err1 != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte("Something went wrong"))
		return
	}
	if resp.ID == 0 {
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte("This record not found in our database"))
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
