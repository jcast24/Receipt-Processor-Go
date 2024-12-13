package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"unicode"

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

func CheckAlphanumeric(r *Receipt) int {
    count := 0
    for _, s := range r.Retailer {
        if unicode.IsLetter(s) || unicode.IsDigit(s) {
            count++
        }
    }
    return count
}

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

	response := struct {
		ID string `json:"id"`
	}{
		ID: newReceipt.ID,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}


func GetPointsById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

    params := mux.Vars(r)
    receiptID := params["id"]

	var foundReceipt *Receipt
    for _, receipt := range receiptStorage {
        if receipt.ID == receiptID {
            foundReceipt = &receipt
            break
        }
    }

    if foundReceipt == nil {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return 
    }

    // if err := json.NewDecoder(r.Body).Decode(&newReceipt); err != nil {
    //     http.Error(w, "Invalid Input", http.StatusBadRequest)
    //     return
    // }

	checkNameResult := CheckAlphanumeric(foundReceipt)
	finalPoints := checkNameResult

	response := struct {
		Points int `json:"points"`
	}{
		Points: finalPoints,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Welcome to my homepage!")
}

func main() {
	r := mux.NewRouter()

	r.HandleFunc("/", homePage).Methods("GET")
	r.HandleFunc("/receipts/process", CreateIdForReceipt).Methods("POST")
	r.HandleFunc("/receipts/{id}/points", GetPointsById).Methods("GET")

	log.Println("Server is starting on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", r))
}
