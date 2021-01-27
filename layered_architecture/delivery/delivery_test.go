package delivery

import (
	"bytes"
	"encoding/json"
	"errors"
	"io/ioutil"
	"layered/architecture/entities"
	"net/http"

	"layered/architecture/service"

	"github.com/gorilla/mux"

	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCustomerGetById(t *testing.T) {
	testcases := []struct {
		id       string
		response entities.Customer
		code     int
	}{
		{"69", entities.Customer{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}, http.StatusOK},
	}

	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/customer/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": testcases[i].id})
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.GetById(w, req)

		var cust entities.Customer

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		//fmt.Println("got ", string(body), " and expected ", v.response)
		_ = json.Unmarshal(body, &cust)

		if !reflect.DeepEqual(cust, v.response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, cust, v.response)
		}

		if w.Code != testcases[i].code {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, testcases[i].code)
		}
	}
}

func TestCustomerGetByName(t *testing.T) {
	testcases := []struct {
		name     string
		response []entities.Customer
		code     int
	}{
		{"bittu%20ray", []entities.Customer{{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}}, http.StatusOK},
		{"", []entities.Customer{{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}}, http.StatusOK},
	}

	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodGet, "/customer?name="+testcases[i].name, nil)
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.GetByName(w, req)

		var cust []entities.Customer

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &cust)

		if !reflect.DeepEqual(cust, v.response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, cust, v.response)
		}

		if w.Code != testcases[i].code {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, testcases[i].code)
		}
	}
}

func TestCustomerCreate(t *testing.T) {
	testcases := []struct {
		inp      []byte
		response entities.Customer
		code     int
	}{
		{[]byte(`{"name":"Pintu","dob":"05/12/2000","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), entities.Customer{ID: 78, Name: "Pintu", DOB: "05-12-2000", Address: entities.Address{ID: 73, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 78}}, http.StatusCreated},
	}

	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodPost, "/customer", bytes.NewBuffer(testcases[i].inp))
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PostCustomer(w, req)

		var cust entities.Customer

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &cust)

		if !reflect.DeepEqual(cust, v.response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, cust, v.response)
		}

		if w.Code != testcases[i].code {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, testcases[i].code)
		}
	}
}

func TestCustomerPut(t *testing.T) {
	testcases := []struct {
		id       string
		inp      []byte
		response entities.Customer
		code     int
	}{
		{"78", []byte(`{"name":"Pintu","address":{"streetName":"Bikramganj","city":"Sasaram","state":"Bihar"}}`), entities.Customer{ID: 78, Name: "Pintu", DOB: "05-12-2000", Address: entities.Address{ID: 73, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 78}}, http.StatusOK},
	}

	for i, v := range testcases {
		req := httptest.NewRequest(http.MethodPut, "/customer/", bytes.NewBuffer(testcases[i].inp))
		req = mux.SetURLVars(req, map[string]string{"id": testcases[i].id})
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.PutCustomer(w, req)

		var cust entities.Customer

		resp := w.Result()
		body, _ := ioutil.ReadAll(resp.Body)
		_ = json.Unmarshal(body, &cust)

		if !reflect.DeepEqual(cust, v.response) {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, cust, v.response)
		}

		if w.Code != testcases[i].code {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, testcases[i].code)
		}
	}
}

func TestCustomerDelete(t *testing.T) {
	testcases := []struct {
		id       string
		response entities.Customer
		code     int
	}{
		{"69", entities.Customer{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}, http.StatusNoContent},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodDelete, "/customer/", nil)
		req = mux.SetURLVars(req, map[string]string{"id": testcases[i].id})
		w := httptest.NewRecorder()

		a := New(mockDatastore{})

		a.DeleteCustomer(w, req)

		if w.Code != testcases[i].code {
			t.Errorf("[TEST%d]Failed. Got %v\tExpected %v\n", i+1, w.Code, testcases[i].code)
		}
	}
}

type mockDatastore struct{}

func (c mockDatastore) GetByID(id int) (entities.Customer, error) {
	data := entities.Customer{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}
	if id == 69 {
		return data, nil
	}
	return entities.Customer{}, errors.New("Something went wrong")
}

func (c mockDatastore) GetByName(name string) ([]entities.Customer, error) {
	data := []entities.Customer{{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}}
	if name == "bittu ray" {
		return data, nil
	}
	if name == "" {
		return data, nil
	}
	return []entities.Customer{}, errors.New("something went wrong")
}

func (c mockDatastore) GetAll() ([]entities.Customer, error) {
	return []entities.Customer{}, nil
}

func (c mockDatastore) CreateCustomer(cust entities.Customer) (entities.Customer, error) {
	if timestamp := service.DateSubstract(cust.DOB); timestamp/(3600*24*12*30) < 18 {
		return entities.Customer{}, errors.New("You are below 18, so you are not allowed to be our customer")
	}
	data := entities.Customer{ID: 78, Name: "Pintu", DOB: "05-12-2000", Address: entities.Address{ID: 73, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 78}}
	return data, nil
}

func (c mockDatastore) UpdateCustomer(id int, customer entities.Customer) (entities.Customer, error) {
	data := entities.Customer{ID: 78, Name: "Pintu", DOB: "05-12-2000", Address: entities.Address{ID: 73, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 78}}
	if id == 78 {
		return data, nil
	}
	return entities.Customer{}, errors.New("Something went wrong")
}

func (c mockDatastore) DeleteCustomer(id int) (entities.Customer, error) {
	if id == 69 {
		return entities.Customer{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}, nil
	}
	return entities.Customer{}, errors.New("some error")
}
