package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

type Product struct {
	Id string
	Name string
	Quantity int
	Price float64
}

var Products []Product

func homepage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to homepage")
	fmt.Println("Endpoint hit: homepage")
}

func returnAllProducts(w http.ResponseWriter, r *http.Request) {
	log.Println("Endpoint hit: returnAllProducts")
	json.NewEncoder(w).Encode(Products)
}

func getProduct(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL.Path)

	vars := mux.Vars(r)

	key := vars["id"]
	for _, product := range Products {
		if string(product.Id) == key {
			json.NewEncoder(w).Encode(product)
		}
	}
}

func handleRequests() {
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/products", returnAllProducts)
	myRouter.HandleFunc("/product/{id}", getProduct)
	myRouter.HandleFunc("/", homepage)
	http.ListenAndServe(":10000", myRouter)
}

func checkError(e error) {
	if e != nil {
		log.Fatalln(e)
	}
}

type Data struct {
	id int
	name string
}

func main() {

	connectionString := fmt.Sprintf("%v:%v@tcp(127.0.0.1:3306)/%v", DbUser, DbPassword, DBName)
	db, err := sql.Open("mysql", connectionString)
	checkError(err)
	defer db.Close()

	rows, err := db.Query("SELECT * from data")
	checkError(err)
	for rows.Next() {
		var data Data
		err := rows.Scan(&data.id, &data.name)
		checkError(err)
		fmt.Println(data)
	}

	Products = []Product{
		{Id: "1", Name: "Chair", Quantity: 100, Price: 100.00},
		{Id: "2", Name: "Desk", Quantity: 200, Price: 200.00},

	}

	handleRequests()

}