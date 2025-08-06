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

	// Example: Get account information

	Processor(client)

}

func Processor(client *binance_connector.Client) {
	for {
		/*
		if RSIstrat(client, "ETHUSDC", "1m", 14, -25) {
			trade := Trade{
				Asset:  "ETHUSDC",
				Amount: 0.002,
			}

			err := trade.Buy(client)
			if err != nil {
				fmt.Println(err)
			}
			log.Println(trade)
			activeTrade = append(activeTrade, trade)
		}

		time.Sleep(1 * time.Minute)
		// check for active Trades

		if len(activeTrade) > 0 {

			// check for sell condtions
			if RSIstrat(client, "ETHUSDC", "1m", 14, 75) {

				err := activeTrade[0].Sell(client)
				SaveTrade(activeTrade[0])
				if err != nil {
					fmt.Println(err)
				}
			}

		}

	}
		*/
}
