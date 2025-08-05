package main

import (
	"log"
	"os"
	"time"

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

	Processor(client)

}

func Processor(client *binance_connector.Client) {

	trade := Trade{}
	trade.Amount = 0.1
	trade.Asset = "ETHUSDC"

	var buyCondition = false
	var sellCondition = false
	for !buyCondition {
		time.Sleep(1 * time.Second)
		buyCondition = true
	}
	err := trade.Buy(client)
	for !sellCondition {
		time.Sleep(2 * time.Second)
		sellCondition = true
	}
	err = trade.Sell(client)

	if err != nil {
		println(err)
	}

	err, _ = GetAssetBalance(client, "ETH")

	if err != nil {
		println(err)
	}

}
