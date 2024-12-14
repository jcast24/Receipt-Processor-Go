package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
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
	isMultipleOfQuarter := math.Mod(floatValue, 0.25) == 0

	if isWholeNumber && isMultipleOfQuarter {
		return 75
	} else if isWholeNumber {
		return 50
	} else if isMultipleOfQuarter {
		return 25
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

	for _, values := range r.Items {
		description := values.ShortDescription
		trimDescription := strings.Trim(description, " ")
		descriptionLength := len(trimDescription)
		var convertedLength float64 = float64(descriptionLength)

		price := values.Price

		// convert price from string to float
		floatPrice, err := strconv.ParseFloat(price, 64)
		if err != nil {
			fmt.Println("Error: ", err)
		}

		var multiplier = 0.25

		if math.Mod(convertedLength, 3) == 0 {
			result = floatPrice * multiplier
			roundedResult = math.Round(result)
			roundedResultUp = int(math.Ceil(roundedResult))
		}
		final += roundedResultUp
	}
	return final
}

func CheckDate(r *Receipt) int {
	var itemPurchaseDate string = r.PurchaseDate
	layout := "2006-01-02"

	date, err := time.Parse(layout, itemPurchaseDate)

	if err != nil {
		fmt.Println("Error parsing date:", err)
	}

	day := date.Day()

	if day%2 == 1 {
		return 6
	}
	return 0
}

func CheckTime(r *Receipt) int {
    var itemPurchaseTime string = r.PurchaseTime
    layout := "15:04"

    startTime, err := time.Parse(layout, "14:00")
    if err != nil {
        fmt.Println("Error parsing: ", err)
    }

    endTime, err := time.Parse(layout, "16:00")
    if err != nil {
        fmt.Println("Error parsing: ", err)
    }

    convertPurchaseTime, err := time.Parse(layout, itemPurchaseTime)
    if err != nil {
        fmt.Println("Error parsing: ", err)
    }

    if !convertPurchaseTime.Before(startTime) && !convertPurchaseTime.After(endTime){
        return 10
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
	checkItemsResult := CheckItemsCount(foundReceipt)
    checkDescriptionResult := CheckDescription(foundReceipt)
    checkDateResult := CheckDate(foundReceipt)
    checkTimeResult := CheckTime(foundReceipt)

	finalPoints := checkNameResult + checkTotalResult + checkItemsResult + checkDescriptionResult + checkDateResult + checkTimeResult

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
