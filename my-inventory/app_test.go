package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	err := a.Initialize(DBUser, DBPassword, "test")
	createTable()
	if err != nil {
		log.Fatal("Error occured while initializing the test database")
	}
	m.Run()
}

func createTable() {
	createTableQuery := `CREATE TABLE IF NOT EXISTS products ( id int NOT NULL AUTO_INCREMENT, name varchar (255) NOT NULL, quantity int, price float (10,7), PRIMARY KEY (id));`
	_, err := a.DB.Exec(createTableQuery)
	if err != nil {
		log.Fatal(err)
	}

}

func clearTable() {
	a.DB.Exec("DELETE from products")
	a.DB.Exec("ALTER table products AUTO_INCREMENT=1")
}

func addProduct(name string, quantity int, price float64) {
	query := fmt.Sprintf("INSERT into products(name, quantity, price) VALUES('%v', %v, %v)", name, quantity, price)
	_, err := a.DB.Exec(query)
	if err != nil {
		log.Fatal(err)
	}

}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProduct("keyboard", 100, 500)

	request, _ := http.NewRequest("GET", "/product/1", nil)

	response := sendRequest(request)
	checkStatusCode(t, http.StatusOK, response.Code)
}

func checkStatusCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected %v, received %v", expected, actual)
	}
}

func sendRequest(request *http.Request) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	a.Router.ServeHTTP(recorder, request)
	return recorder
}

func TestCreateProduct(t *testing.T) {
	clearTable()
	product := []byte(`{"name": "chair", "quantity": 1, "price": 100}`)
	req, _ := http.NewRequest("POST", "/product", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response := sendRequest(req)
	checkStatusCode(t, http.StatusCreated, response.Code)

	var m map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &m)
	
	if m["name"] != "chair" {
		t.Errorf("Expected name: %v, Got: %v", "chair", m["name"])
	}

	if m["quantity"] != 1.0 {
	t.Errorf(  "Expected quantity: %v, Got: %v", 1.0, m["quantity"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 100)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("DELETE", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)

	req, _ = http.NewRequest("GET", "/product/1", nil)
	response = sendRequest(req)
	checkStatusCode(t, http.StatusNotFound, response.Code)
		
}

func TestUpdateProduct(t *testing.T) {
	clearTable()
	addProduct("connector", 10, 100)

	req, _ := http.NewRequest("GET", "/product/1", nil)
	response := sendRequest(req)
	checkStatusCode(t, http.StatusOK, response.Code)
	var oldValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &oldValue)

	product := []byte(`{"name": "connector", "quantity": 1, "price": 5}`)
	req, _ = http.NewRequest("PUT", "/product/1", bytes.NewBuffer(product))
	req.Header.Set("Content-Type", "application/json")
	response = sendRequest(req)

	checkStatusCode(t, http.StatusOK, response.Code)

	var newValue map[string]interface{}
	json.Unmarshal(response.Body.Bytes(), &newValue)

	if oldValue["id"] != newValue["id"] {
		t.Errorf("Expected id: %v, got: %v", oldValue["id"], newValue["id"])
	}

	if newValue["quantity"] != 1.0 {
		t.Errorf("Expected quantity: %v, got: %v", 1, newValue["quantity"])
	}

	if newValue["price"] != 5.0 {
		t.Errorf("Expected price: %v, got: %v", 5, newValue["price"])
	}

}

// check product not exist in the db
// try an getAll products
// try malformed payloads
// try deleting a product that doesnt exist
	