package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	// Initialize Binance client
	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")
	_ = client

	var testStrategy = Strategy{
		Asset:    "ETHUSDC",
		Interval: "1h",
		Filters:  []Filter{},
		Main: Signal{
			Name: "SMA",
			Type: "Moving Average",
			Params: map[string]int{
				"short": 9,
				"long":  20,
			},
		},
	}
	result := testStrategy.StrategyTester(client)

	jsonBytes, err := json.Marshal(result)

	file, err := os.OpenFile("startReport.json", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)

	if err != nil {
		fmt.Println(err)
	}
	file.Write(jsonBytes)
	if err != nil {
		fmt.Println(err)
	}

	defer file.Close()

}
