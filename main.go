// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/shopspring/decimal"
	"github.com/prometheus/client_golang/prometheus/promhttp"

)

var db *gorm.DB
var err error

// Product is a representation of a product
type Product struct {
	ID    int             `form:"id" json:"id"`
	Code  string          `form:"code" json:"code"`
	Name  string          `form:"name" json:"name"`
	Price decimal.Decimal `form:"price" json:"price" sql:"type:decimal(16,2);"`
}

// Result is an array of product
type Result struct {
	Code    int         `json:"code"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

func main() {
	
	db, err = gorm.Open("mysql", "root:r00tpa55@(db:3306)/learning?charset=utf8&parseTime=True")

	if err !=nil {
		log.Println("connection failed",err)
	} else {
		log.Println("connection established")
	}
	db.AutoMigrate(&Product{})
	handleRequests()
	
}
func handleRequests() {
	log.Println("Start the server at http://0.0.0.0:9999")

	// http.Handle("/metrics", promhttp.Handler())
    // http.ListenAndServe(":2112", nil)

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/healthz", homePage)
	myRouter.HandleFunc("/products", createProduct).Methods("POST")
	myRouter.HandleFunc("/products", getProducts).Methods("GET")
	myRouter.Handle("/metrics", promhttp.Handler())

	log.Fatal(http.ListenAndServe(":9999", myRouter))
}


func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "APP IS RUNNING PROPERLY!")
}
func createProduct(w http.ResponseWriter, r *http.Request) {
	payloads, _ := ioutil.ReadAll(r.Body)
	fmt.Fprintf(w, "create product")
	var product Product
	json.Unmarshal(payloads, &product)

	db.Create(&product)

	res := Result{Code: 200, Data: product, Message: "Success create product"}
	result, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(result)
}

func getProducts(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Endpoint hit: get products")

	products := []Product{}
	db.Find(&products)

	res := Result{Code: 200, Data: products, Message: "Success get products"}
	results, err := json.Marshal(res)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(results)
}

