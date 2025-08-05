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

	// Get API credentials from environment variables
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	// Initialize Binance client
	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")

	// Example: Get account information

	if err != nil {
		log.Fatalf("Error getting account info: %v", err)
	}

	trade := Trade{}
	trade.Amount = 0.1
	trade.Asset = "ETHUSDC"
	if err != nil {
		println(err)
	}
	err = trade.Buy(client)

	if err != nil {
		println(err)
	}

	err, balance := GetAssetBalance(client, "ETH")

	if err != nil {
		println(err)
	}

	fmt.Printf("Decimal (%%f): %f\n", balance)

	// Example: Get current price for a trading pair
	_ = "BTCUSDC"

}
