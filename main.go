package main

import (
	"fmt"
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
	_ = client

	balance, _ := GetAssetBalance(client, "USDC")
	fmt.Println(balance)
	Processor(client)

}

func Processor(client *binance_connector.Client) {

	strategies := []*Strategy{
		{
			Asset:         "HBARUSDC",
			Amount:        1000,
			BuyCondition:  RSIbuyCondition14,
			SellCondition: RSIsellCondition14,
		},
		{
			Asset:         "ETHUSDC",
			Amount:        0.1,
			BuyCondition:  RSIbuyCondition14,
			SellCondition: RSIsellCondition14,
		},
	}

	totalTrade := len(strategies)

	for totalTrade > 0 {

		for _, s := range strategies {
			if !s.TradeInProgress && s.BuyCondition(client, s.Asset, "1m") {
				err := s.Buy(client)
				if err != nil {
					fmt.Println(err)
				}

			}

			time.Sleep(50 * time.Second)

			if s.TradeInProgress && s.BuyCondition(client, s.Asset, "1m") {
				err := s.Sell(client)
				if err != nil {
					fmt.Println(err)
				}

			}

		}
	}
}
