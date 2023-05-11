package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type Receipt struct {
	Retailer     string    `json:"retailer"`
	PurchaseDate string    `json:"purchaseDate"`
	PurchaseTime string    `json:"purchaseTime"`
	Items        []ReceiptItem `json:"items"`
	Total        string    `json:"total"`
}

type ReceiptItem struct {
	ShortDescription string `json:"shortDescription"`
	Price            string `json:"price"`
}

type PointsResponse struct {
	Points int `json:"points"`
}

func main() {
	http.HandleFunc("/receipts/process", processReceiptsHandler)
	http.HandleFunc("/receipts/", getPointsHandler)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func processReceiptsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	var receipt Receipt
	err := json.NewDecoder(r.Body).Decode(&receipt)
	if err != nil {
		http.Error(w, "Invalid request payload.", http.StatusBadRequest)
		return
	}

	id := uuid.New().String()
	response := map[string]string{"id": id}

	// Save receipt and id to some in-memory data store for later retrieval

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func getPointsHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method.", http.StatusMethodNotAllowed)
		return
	}

	path := strings.TrimPrefix(r.URL.Path, "/receipts/")
	id := strings.TrimSuffix(path, "/points")

	// Retrieve receipt from in-memory data store using id

	var receipt Receipt

	// Calculate points based on receipt data
	points := 0

	// Rule 1: One point for every alphanumeric character in the retailer name.
	points += len(strings.ReplaceAll(receipt.Retailer, " ", ""))

	// Rule 2: 50 points if the total is a round dollar amount with no cents.
	totalFloat, err := strconv.ParseFloat(receipt.Total, 64)
	if err == nil {
		if totalFloat == float64(int(totalFloat)) {
			points += 50
		}
	}

	// Rule 3: 25 points if the total is a multiple of 0.25.
	totalMod := totalFloat / 0.25
	if totalMod == float64(int(totalMod)) {
		points += 25
	}

	// Rule 4: 5 points for every two items on the receipt.
	points += (len(receipt.Items) / 2) * 5

	// Rule 5: If the trimmed length of the item description is a multiple of 3, multiply the price by 0.2 and round up to the nearest integer.
	for _, item := range receipt.Items {
		trimmedLen := len(strings.TrimSpace(item.ShortDescription))
		if trimmedLen%3 == 0 {
			priceFloat, err := strconv.ParseFloat(item.Price, 64)
			if err == nil {
				pricePoints := int(priceFloat * 0.2)
				if priceFloat*0.2-float64(pricePoints) > 0 {
					pricePoints++
				}
				points += pricePoints
			}
		}
	}

	// Rule 6: 6 points if the day in the purchase date is odd.
	purchaseDate, err := time.Parse
