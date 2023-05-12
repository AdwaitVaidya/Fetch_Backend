package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"testing"
)

func insertData() [3]string {
	var inserts = []Receipt{
		{Retailer: "Target",
			PurchaseDate: "2022-01-02",
			PurchaseTime: "13:13",
			Items: []Item{
				{ShortDescription: "Pepsi - 12-oz", Price: 1.25},
			},
			Total: 1.25},
		{Retailer: "Target",
			PurchaseDate: "2022-01-01",
			PurchaseTime: "13:01",
			Items: []Item{
				{
					ShortDescription: "Mountain Dew 12PK",
					Price:            6.49,
				},
				{
					ShortDescription: "Emils Cheese Pizza",
					Price:            12.25,
				},
				{
					ShortDescription: "Knorr Creamy Chicken",
					Price:            1.26,
				},
				{
					ShortDescription: "Doritos Nacho Cheese",
					Price:            3.35,
				},
				{
					ShortDescription: "   Klarbrunn 12-PK 12 FL OZ  ",
					Price:            12.00,
				},
			},
			Total: 35.35},
		{Retailer: "M&M Corner Market",
			PurchaseDate: "2022-03-20",
			PurchaseTime: "14:33",
			Items: []Item{
				{ShortDescription: "Gatorade", Price: 2.25},
				{ShortDescription: "Gatorade", Price: 2.25},
				{ShortDescription: "Gatorade", Price: 2.25},
				{ShortDescription: "Gatorade", Price: 2.25},
			},
			Total: 9.00},
	}
	var arr [3]string
	for index := range inserts {
		payloadBytes, err := json.Marshal(inserts[index])
		if err != nil {
			fmt.Println("Error marshalling payload:", err)
		}
		resp, err := http.Post("http://localhost:8080/receipts/process", "application/json", bytes.NewReader(payloadBytes))
		if err != nil {
			fmt.Println("Error sending POST request:", err)
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
		}
		sb := string(body)
		var data map[string]interface{}
		err = json.Unmarshal([]byte(sb), &data)
		if err != nil {
			fmt.Println("Error parsing JSON:", err)
		}
		id, _ := data["id"].(string)
		arr[index] = id
	}
	return arr
}

func TestGetPoints(t *testing.T) {
	var arr [3]string = insertData()
	//println(id)
	name := "TestGetPoints"
	var res [3]string = [3]string{"31", "28", "100"}
	for index := range arr {
		t.Run(name, func(t *testing.T) {
			resp, err := http.Get("http://localhost:8080/receipts/" + arr[index] + "/points")
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
			sb := string(body)
			var jsonMap map[string]int32
			json.Unmarshal([]byte(sb), &jsonMap)
			var twt int32 = jsonMap["points"]

			//println(val == sb)
			st := res[index]
			num, _ := strconv.Atoi(st)
			if int32(num) != twt {
				t.Errorf("got %q, want %q", num, twt)
			}
		})
	}

}
