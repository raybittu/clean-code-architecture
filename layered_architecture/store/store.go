package store

import "C"
import (
	"database/sql"
	"layered/architecture/entities"
	"layered/architecture/errors"
)

type CustomerStore struct {
	db *sql.DB
}

func (c CustomerStore) Close() {
	c.db.Close()
}

func New() Customer {
	var db, dbErr = sql.Open("mysql", "root:1118209@/Customer_service")
	if dbErr != nil {
		panic(dbErr)
	}
	//defer db.Close()

	dbErr = db.Ping()
	if dbErr != nil {
		panic(dbErr.Error()) // proper error handling instead of panic in your app
	}
	return CustomerStore{db: db}
}

func (c CustomerStore) GetByID(id int) (entities.Customer, error) {
	query := "select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=? order by cust.id, addrs.id"

	var DataInterface []interface{}
	DataInterface = append(DataInterface, id)
	rows, err := c.db.Query(query, DataInterface...)
	if err != nil {
		return entities.Customer{}, errors.ErrDBQuery
	}

	var cust entities.Customer

	for rows.Next() {
		_ = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}

	return cust, nil
}

func (c CustomerStore) GetByName(name string) ([]entities.Customer, error) {
	query := "select * from cust inner join addrs on cust.id=addrs.cus_id order by cust.id, addrs.id"
	var data []interface{}
	if name != "" {
		query = "select * from cust inner join addrs on cust.id=addrs.cus_id where cust.name=? order by cust.id, addrs.id"
		data = append(data, name)
	}

	rows, err := c.db.Query(query, data...)

	if err != nil {
		return []entities.Customer(nil), errors.ErrDBQuery
	}
	var customer []entities.Customer

	for rows.Next() {
		var cust entities.Customer
		_ = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
		customer = append(customer, cust)
	}

	return customer, nil
}

func (c CustomerStore) Create(customer entities.Customer) (entities.Customer, error) {
	query := "insert into cust (name, dob) values (?, ?)"
	var CustomerData []interface{}
	CustomerData = append(CustomerData, customer.Name)
	CustomerData = append(CustomerData, customer.DOB)

	rows, err := c.db.Exec(query, CustomerData...)
	if err != nil {
		return entities.Customer{}, errors.ErrDBExec
	}

	id, err1 := rows.LastInsertId()

	if err1 != nil {
		return entities.Customer{}, errors.ErrDBExec
	} else {
		query = "insert into addrs (streetname, city, state, cus_id) values (?, ?, ?, ?)"
		var AddressData []interface{}
		AddressData = append(AddressData, customer.Address.StreetName)
		AddressData = append(AddressData, customer.Address.City)
		AddressData = append(AddressData, customer.Address.State)
		AddressData = append(AddressData, id)
		rows, err = c.db.Exec(query, AddressData...)
		if err == nil {
			customer.ID = int(id)
			addressId, _ := rows.LastInsertId()
			customer.Address.ID = int(addressId)
			customer.Address.CusId = int(id)
			return customer, nil
		}
		return entities.Customer{}, errors.ErrDBExec
	}
}

func (c CustomerStore) Update(id int, customer entities.Customer) (entities.Customer, error) {

	if customer.Name != "" {
		query := "update cust set name=? where id=?"
		_, err := c.db.Exec(query, customer.Name, id)
		if err != nil {
			return entities.Customer{}, errors.ErrDBExec
		}

	}
	var data []interface{}
	query := "update addrs set "
	if customer.Address.State != "" {
		query += "state = ? ,"
		data = append(data, customer.Address.State)
	}
	if customer.Address.City != "" {
		query += "city = ? ,"
		data = append(data, customer.Address.City)
	}
	if customer.Address.StreetName != "" {
		query += "streetname = ? ,"
		data = append(data, customer.Address.StreetName)
	}
	query = query[:len(query)-1]
	query += "where cus_id = ?"
	data = append(data, id)
	_, err := c.db.Exec(query, data...)

	if err != nil {
		return entities.Customer{}, errors.ErrDBExec
	}
	query = "select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=?"
	rows, _ := c.db.Query(query, id)
	var cust entities.Customer
	for rows.Next() {
		_ = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}
	return cust, nil
}

func (c CustomerStore) Delete(id int) (entities.Customer, error) {
	query := "select * from cust inner join addrs on addrs.cus_id=cust.id and cust.id=? order by cust.id, addrs.id"
	rows, err := c.db.Query(query, id)
	if err != nil {
		return entities.Customer{}, errors.ErrDBQuery
	}
	var cust entities.Customer
	for rows.Next() {
		_ = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}

	query = "delete from cust where id=?"
	_, err = c.db.Exec(query, id)
	if err != nil {
		return entities.Customer{}, errors.ErrDBExec
	}
	return cust, nil
}
