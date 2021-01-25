package service

import (
	"layered/architecture/entities"
)

type Customer interface {
	GetByID(id int) (entities.Customer, error)
	GetAll() ([]entities.Customer, error)
	GetByName(name string) ([]entities.Customer, error)
	CreateCustomer(c entities.Customer) (entities.Customer, error)
	UpdateCustomer(id int, c entities.Customer) (entities.Customer, error)
	DeleteCustomer(id int) (entities.Customer, error)
}
