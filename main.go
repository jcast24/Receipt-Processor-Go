package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

// full recipe
type Receipt struct {
	ID           string `json:"id"`
	Retailer     string `json:"retailer"`
	PurchaseDate string `json:"purchaseDate"`
	PurchaseTime string `json:"purchaseTime"`
	Items        []Item `json:"items"`
	Total        string `json:"total"`
}

// represents individual items part of the receipt
type Item struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

// In memory storage
var receiptStorage []Receipt

func CreateIdForReceipt(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newReceipt Receipt

	// Decode JSON payload
	err := json.NewDecoder(r.Body).Decode(&newReceipt)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// generate UUID
	newReceipt.ID = uuid.New().String()

	receiptStorage = append(receiptStorage, newReceipt)

	w.WriteHeader(http.StatusCreated)
	// w.Write([]byte(newReceipt.ID))
    
    response := struct {
        ID string `json:"id"`
    } {
        ID: newReceipt.ID,
    }

    json.NewEncoder(w).Encode(response)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my homepage!")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/receipts/process", CreateIdForReceipt).Methods("POST")

	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
