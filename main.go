package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
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

func CheckTotal(r *Receipt) int {
    total := r.Total

    floatValue, err := strconv.ParseFloat(total, 64)
    if err != nil {
        fmt.Println("Error: ", err)
    }

    isWholeNumber := math.Mod(floatValue, 1) == 0
    isMultipleOfQuarter :=  math.Mod(floatValue, 0.25) == 0

    if (isWholeNumber && isMultipleOfQuarter) {
        return 75
    } else if (isWholeNumber) {
        return 50;
    } else if (isMultipleOfQuarter) {
        return 25;
    }
    return 0
}

func CheckItemsCount(r *Receipt) int {
    count := (len(r.Items) / 2) * 5
    return count
}

func CheckDescription(r *Receipt) int {
    var result = 0.0
    var final = 0
    var roundedResult = 0.0
    var roundedResultUp = 0

    for _,values := range r.Items {
        description := values.ShortDescription 
        trimDescription := strings.Trim(description, " ")
        

        price := values.Price
        
        // convert price from string to float
        floatPrice, err := strconv.ParseFloat(price, 64)
        if err != nil {
            fmt.Println("Error: ", err)
        }
    
        var multiplier = 0.25

        if (math.Mod(len(trimDescription), 3) == 0) {
            result = floatPrice * multiplier
            roundedResult = math.Round(result)
        }

    }


    return 0
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
    checkTotalResult := CheckTotal(foundReceipt)
    checkItemsResult := CheckTotal(foundReceipt)

	finalPoints := checkNameResult + checkTotalResult + checkItemsResult

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
