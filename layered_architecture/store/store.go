package store

import "C"
import (
	"database/sql"
	"layered/architecture/entities"
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
	rows, err := c.db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=? order by cust.id, addrs.id", id)
	if err != nil {
		return entities.Customer{}, err
	}

	var cust entities.Customer

	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
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
		return []entities.Customer(nil), err
	}
	var customer []entities.Customer

	for rows.Next() {
		var cust entities.Customer
		err = rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
		customer = append(customer, cust)
	}

	return customer, nil
}

func (c CustomerStore) Create(customer entities.Customer) (entities.Customer, error) {

	rows, err := c.db.Exec("insert into cust (name, dob) values (?, ?)", customer.Name, customer.DOB)
	if err != nil {
		return entities.Customer{}, err
	}

	id, err1 := rows.LastInsertId()

	if err1 != nil {
		return entities.Customer{}, err
	} else {
		rows, err = c.db.Exec("insert into addrs (streetname, city, state, cus_id) values (?, ?, ?, ?)", customer.Address.StreetName, customer.Address.City, customer.Address.State, id)
		if err == nil {
			customer.ID = int(id)
			addressId, _ := rows.LastInsertId()
			customer.Address.ID = int(addressId)
			customer.Address.CusId = int(id)
			return customer, nil
		}
		return entities.Customer{}, err
	}
}

func (c CustomerStore) Update(id int, customer entities.Customer) (entities.Customer, error) {
	if customer.Name != "" {
		_, err := c.db.Exec("update cust set name=? where id=?", customer.Name, id)
		if err != nil {
			return entities.Customer{}, nil
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
		return entities.Customer{}, err
	}

	rows, _ := c.db.Query("select * from cust inner join addrs on cust.id=addrs.cus_id and cust.id=?", id)
	var cust entities.Customer
	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}
	return cust, nil
}

func (c CustomerStore) Delete(id int) (entities.Customer, error) {
	rows, err := c.db.Query("select * from cust inner join addrs on addrs.cus_id=cust.id and cust.id=? order by cust.id, addrs.id", id)
	if err != nil {
		return entities.Customer{}, err
	}
	var cust entities.Customer
	for rows.Next() {
		rows.Scan(&cust.ID, &cust.Name, &cust.DOB, &cust.Address.ID, &cust.Address.StreetName, &cust.Address.City, &cust.Address.State, &cust.Address.CusId)
	}
	_, err = c.db.Exec("delete from cust where id=?", id)
	if err != nil {
		return entities.Customer{}, err
	}
	return cust, nil
}
