package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Receipt struct {
	Retailer     string  `json:"retailer"`
	PurchaseDate string  `json:"purchaseDate"`
	PurchaseTime string  `json:"purchaseTime"`
	Items        []Item  `json:"items"`
	Total        float64 `json:"total"`
}

// Struct to represent an item on a receipt
type Item struct {
	ShortDescription string  `json:"shortDescription"`
	Price            float64 `json:"price"`
}

// Struct to represent the response containing the ID of the receipt
type ResponseID struct {
	ID string `json:"id"`
}

// Struct to represent the response containing the number of points awarded
type ResponsePoints struct {
	Points int `json:"points"`
}

// store points with specific uid (uid)->(points)
var myMap map[string]int

/*
Initializa the myMap and attach functions to endpoints.
*/
func main() {
	router := mux.NewRouter()
	myMap = make(map[string]int)
	router.HandleFunc("/receipts/process", ProcessReceipt).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", GetPoints).Methods("GET")
	fmt.Println("Server started on port 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}

/*
Process the receipt, calculate points, uid and set value in map. Returns a ResposeID with
specific uuid.
*/
func ProcessReceipt(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	u := generateID()
	myMap[u] = CalculatePoints(receipt)
	var data ResponseID
	data.ID = u

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(data)
}

/*
Function to calculate points from a given receipt.
*/
func CalculatePoints(receipt Receipt) int {
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points += len(strings.ReplaceAll(receipt.Retailer, " ", ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	total := fmt.Sprintf("%.2f", receipt.Total)
	if strings.HasSuffix(total, ".00") {
		points += 50
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	if math.Mod(receipt.Total, 0.25) == 0 {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += len(receipt.Items) / 2 * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen > 0 && trimmedLen%3 == 0 {
			points += int(math.Ceil(item.Price * 0.2))
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse("2006-01-02", receipt.PurchaseDate)
	if err == nil && purchaseDate.Day()%2 == 1 {
		points += 6
	}

	// Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
	purchaseTime, err := time.Parse("15:04", receipt.PurchaseTime)
	if err == nil {
		if purchaseTime.After(time.Date(1, 1, 1, 14, 0, 0, 0, time.UTC)) && purchaseTime.Before(time.Date(1, 1, 1, 16, 0, 0, 0, time.UTC)) {
			points += 10
		}
	}

	return points
}

/*
Generate point and
*/
func GetPoints(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	receiptID := params["id"]
	// Look up the receipt by ID and calculate the points awarded
	// ...
	var points ResponsePoints // Replace with the actual points awarded
	ok := false
	points.Points, ok = myMap[receiptID]

	if !ok {
		http.Error(w, "key not in map", http.StatusInternalServerError)
		return
	}
	// Create the response object

	// Encode the response object as JSON and send it in the response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(points)
}

/*
Generate and return uuid.
*/
func generateID() string {
	return uuid.New().String()
}
