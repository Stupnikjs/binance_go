package main

import (
	"fmt"
	"log"
	"os"
	"strings"

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
	// Processor(client)

	for {

		cmd := Prompt()
		arr := strings.Split(cmd, " ")
		if arr[0] == "GET" {
			balance, err := GetAssetBalance(client, arr[1])
			if err != nil {
				fmt.Println(err)
			} else {
				fmt.Println(balance)
			}
		}
		if arr[0] == "TRADE" {
			Processor(client)
		}

	}

}

func Processor(client *binance_connector.Client) {

	robot := Bot{
		client: client,
		strategies: []*Strategy{
			{
				Asset:         "BTCUSDC",
				Amount:        0.001,
				BuyCondition:  RSIbuyCondition14,
				SellCondition: RSIsellCondition14,
			},
			{
				Asset:         "ETHUSDC",
				Amount:        0.001,
				BuyCondition:  RSIbuyCondition14,
				SellCondition: RSIsellCondition14,
			},
			{
				Asset:         "HBARUSDC",
				Amount:        1000,
				BuyCondition:  RSIbuyCondition14,
				SellCondition: RSIsellCondition14,
			},
			{
				Asset:         "LINKUSDC",
				Amount:        40,
				BuyCondition:  RSIbuyCondition14,
				SellCondition: RSIsellCondition14,
			},
		},
	}
	robot.Run()

}
