package store

import (
	"layered/architecture/entities"
	"reflect"
	"testing"

	_ "github.com/go-sql-driver/mysql"
)

func TestDatastore(t *testing.T) {
	a := New()
	testCustomerGetByID(t, a)
	testCustomerGetByName(t, a)
	testCustomerCreate(t, a)
	testCustomerUpdate(t, a)
	testCustomerDelete(t, a)

}

func testCustomerGetByID(t *testing.T, db Customer) {
	testcases := []struct {
		req  int
		resp entities.Customer
	}{
		{69, entities.Customer{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}},
		{70, entities.Customer{ID: 70, Name: "Bittu", DOB: "10/10/2000", Address: entities.Address{ID: 65, StreetName: "Itadhiyan", City: "Bikramganj", State: "Bihar", CusId: 70}}},
	}
	for i := range testcases {
		res, err := db.GetByID(testcases[i].req)

		if err != nil {
			t.Errorf("[TEST%d] FAILED!, got error %v\n", i+1, err)
		}

		if !reflect.DeepEqual(res, testcases[i].resp) {
			t.Errorf("[TEST%d] FAILED!, Wanted %v\t got %v\n", i+1, testcases[i].resp, res)
		}
	}
}

func testCustomerGetByName(t *testing.T, db Customer) {
	testcases := []struct {
		req  string
		resp []entities.Customer
	}{
		{"bittu ray", []entities.Customer{{ID: 69, Name: "bittu ray", DOB: "10/10/2000", Address: entities.Address{ID: 64, StreetName: "HSR", City: "bangaluru", State: "karnataka", CusId: 69}}}},
	}
	for i := range testcases {
		res, err := db.GetByName(testcases[i].req)

		if err != nil {
			t.Errorf("[TEST%d] FAILED!, got error %v\n", i+1, err)
		}

		if !reflect.DeepEqual(res, testcases[i].resp) {
			t.Errorf("[TEST%d] FAILED!, Wanted %v\t got %v\n", i+1, testcases[i].resp, res)
		}
	}
}

func testCustomerCreate(t *testing.T, db Customer) {
	testcases := []struct {
		req  entities.Customer
		resp entities.Customer
	}{
		{entities.Customer{Name: "Kevin", DOB: "15/12/2000", Address: entities.Address{StreetName: "Bikramganj", City: "Sasaram", State: "Bihar"}}, entities.Customer{ID: 75, Name: "Kevin", DOB: "15/12/2000", Address: entities.Address{ID: 70, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 75}}},
	}
	for i := range testcases {
		res, err := db.Create(testcases[i].req)
		if err != nil {
			t.Errorf("[TEST%d] FAILED!, got error %v\n", i+1, err)
		}

		if !reflect.DeepEqual(res, testcases[i].resp) {
			t.Errorf("[TEST%d] FAILED!, Wanted %v\t got %v\n", i+1, testcases[i].resp, res)
		}
	}
}

func testCustomerUpdate(t *testing.T, db Customer) {
	testcases := []struct {
		id   int
		req  entities.Customer
		resp entities.Customer
	}{
		{73, entities.Customer{Name: "Kevin John", DOB: "15/12/2000", Address: entities.Address{StreetName: "Bikramganj", City: "Sasaram", State: "Bihar"}}, entities.Customer{ID: 73, Name: "Kevin John", DOB: "15/12/2000", Address: entities.Address{ID: 68, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 73}}},
	}
	for i := range testcases {
		res, err := db.Update(testcases[i].id, testcases[i].req)
		if err != nil {
			t.Errorf("[TEST%d] FAILED!, got error %v\n", i+1, err)
		}

		if !reflect.DeepEqual(res, testcases[i].resp) {
			t.Errorf("[TEST%d] FAILED!, Wanted %v\t got %v\n", i+1, testcases[i].resp, res)
		}
	}
}

func testCustomerDelete(t *testing.T, db Customer) {
	testcases := []struct {
		req  int
		resp entities.Customer
	}{
		{74, entities.Customer{ID: 74, Name: "Kevin", DOB: "15/12/2000", Address: entities.Address{ID: 69, StreetName: "Bikramganj", City: "Sasaram", State: "Bihar", CusId: 74}}},
	}
	for i := range testcases {
		res, err := db.Delete(testcases[i].req)
		if err != nil {
			t.Errorf("[TEST%d] FAILED!, got error %v\n", i+1, err)
		}

		if !reflect.DeepEqual(res, testcases[i].resp) {
			t.Errorf("[TEST%d] FAILED!, Wanted %v\t got %v\n", i+1, testcases[i].resp, res)
		}
	}
}
