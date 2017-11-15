package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// Item represents shopping item (1)
type Item struct {
	Name        string  `json:"name"`
	Supermarket string  `json:"supermarket"`
	Price       float64 `json:"price"`
}

// ShoppingItems are the items in a shopping list represented as a slice (dynamic array)
var ShoppingItems []Item

// Index returns sample of items
func Index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	item1 := Item{Name: "Egg", Supermarket: "Netto", Price: 20.0}
	item2 := Item{Name: "Beef", Supermarket: "Irma", Price: 80.0}
	if len(ShoppingItems) == 0 {
		ShoppingItems = append(ShoppingItems, item1, item2)
	}

	json.NewEncoder(w).Encode(ShoppingItems)
}

// AddItem adds item to shopping list (2 REST METHODS)
func AddItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	decoder := json.NewDecoder(r.Body)
	var item Item
	err := decoder.Decode(&item)
	if err != nil {
		panic(err)
	}
	ShoppingItems = append(ShoppingItems, item)
	//fmt.printf("Item %s successfully added.", item.Name)
}

// DeleteItem removes shopping item
func DeleteItem(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for i, item := range ShoppingItems {
		if item.Name == params["name"] {
			ShoppingItems = append(ShoppingItems[:i], ShoppingItems[i+1:]...)
			break
		}
	}
}

// DeleteAllItems removes all items in shopping list
func DeleteAllItems(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	ShoppingItems = ShoppingItems[:0]
}

// GetTotalItemPrice returns the total price based on all shopping items
func GetTotalItemPrice(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var totalPrice = 0.0
	for _, item := range ShoppingItems {
		totalPrice += item.Price
	}
	json.NewEncoder(w).Encode(totalPrice)
}

// GetAllItemsFromSupermarket returns items for a specific shopping list
func GetAllItemsFromSupermarket(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	items := make([]Item, 0)
	for _, item := range ShoppingItems {
		if strings.TrimSpace(strings.ToLower(item.Supermarket)) == strings.TrimSpace(strings.ToLower(params["supermarket"])) {
			items = append(items, item)
		}
	}
	json.NewEncoder(w).Encode(items)
}

// Server starts listening for HTTP requests
func Server() {
	if err := http.ListenAndServe("127.0.0.1:8080", nil); err != nil {
		log.Fatal(err)
	}
}

// Route routes requests to handlers using Gorilla
func Route() {
	// Handle all requests with the Gorilla router.
	r := mux.NewRouter()
	http.Handle("/", r)
	r.HandleFunc("/items", Index).Methods("GET")                                  // curl localhost:8080/items
	r.HandleFunc("/add", AddItem).Methods("POST")                                 // curl -H "Content-Type: application/json" -X POST -d '{"name":"t", "supermarket”:t", "price":0}' localhost:8080/add
	r.HandleFunc("/delete/{name}", DeleteItem).Methods("POST")                    // curl -d '' localhost:8080/delete/Beef
	r.HandleFunc("/delete", DeleteAllItems).Methods("POST")                       // curl -d '' localhost:8080/delete
	r.HandleFunc("/get", GetTotalItemPrice).Methods("GET")                        // curl localhost:8080/get
	r.HandleFunc("/get/{supermarket}", GetAllItemsFromSupermarket).Methods("GET") // curl localhost:8080/get/irma
}
