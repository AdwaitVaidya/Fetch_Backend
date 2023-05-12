package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
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

func main() {
	//url := "http://localhost:8080/receipts/process"
	payload := Receipt{
		Retailer:     "Target",
		PurchaseDate: "2022-01-02",
		PurchaseTime: "13:13",
		Total:        1.25,
		Items: []Item{
			{ShortDescription: "Pepsi - 12-oz", Price: 1.25},
		},
	}
	// payload = Receipt{
	// 	Retailer:     "Target",
	// 	PurchaseDate: "2022-01-01",
	// 	PurchaseTime: "13:01",
	// 	Items: []Item{
	// 		{
	// 			ShortDescription: "Mountain Dew 12PK",
	// 			Price:            6.49,
	// 		},
	// 		{
	// 			ShortDescription: "Emils Cheese Pizza",
	// 			Price:            12.25,
	// 		},
	// 		{
	// 			ShortDescription: "Knorr Creamy Chicken",
	// 			Price:            1.26,
	// 		},
	// 		{
	// 			ShortDescription: "Doritos Nacho Cheese",
	// 			Price:            3.35,
	// 		},
	// 		{
	// 			ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
	// 			Price:            12.00,
	// 		},
	// 	},
	// 	Total: 35.35,
	// }

	// payload = Receipt{
	// 	Retailer:     "M&M Corner Market",
	// 	PurchaseDate: "2022-03-20",
	// 	PurchaseTime: "14:33",
	// 	Items: []Item{
	// 		{ShortDescription: "Gatorade", Price: 2.25},
	// 		{ShortDescription: "Gatorade", Price: 2.25},
	// 		{ShortDescription: "Gatorade", Price: 2.25},
	// 		{ShortDescription: "Gatorade", Price: 2.25},
	// 	},
	// 	Total: 9.00}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("Error marshalling payload:", err)
		return
	}
	resp, err := http.Post("http://localhost:8080/receipts/process", "application/json", bytes.NewReader(payloadBytes))
	if err != nil {
		fmt.Println("Error sending POST request:", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}
	print(body)
	// Unmarshal the response to a ResponseID struct
	// var responseID ResponseID
	// if err := json.Unmarshal(body, &responseID); err != nil {
	// 	fmt.Println("Error unmarshalling response body:", err)
	// 	return
	// }
}
