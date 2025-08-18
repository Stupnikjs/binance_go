package main

import (
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

	var apiKey = os.Getenv("API_KEY")
	var secretKey = os.Getenv("SECRET_KEY")

	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")
	_ = client
	// Get API credentials from environment variables

	err = FetchReports(client, []Interval{m5, m15, m30})
	if err != nil {
		fmt.Println(err)
	}

}

// creer des commandes pour backtester
