package service

import (
	"layered/architecture/entities"
	"layered/architecture/errors"
	"layered/architecture/store"
)

type CustomerService struct {
	store store.Customer
}

func New(customer store.Customer) Customer {
	return CustomerService{store: customer}
}

func (c CustomerService) GetByID(id int) (entities.Customer, error) {
	return c.store.GetByID(id)
}

func (c CustomerService) GetByName(name string) ([]entities.Customer, error) {
	return c.store.GetByName(name)
}

func (c CustomerService) GetAll() ([]entities.Customer, error) {
	return c.store.GetByName("")
}

func (c CustomerService) CreateCustomer(cust entities.Customer) (entities.Customer, error) {
	if timestamp := DateSubstract(cust.DOB); timestamp/(3600*24*12*30) < 18 {

		return entities.Customer{}, errors.ErrEligibility
	}
	return c.store.Create(cust)
}

func (c CustomerService) UpdateCustomer(id int, customer entities.Customer) (entities.Customer, error) {
	return c.store.Update(id, customer)
}

func (c CustomerService) DeleteCustomer(id int) (entities.Customer, error) {
	return c.store.Delete(id)
}
