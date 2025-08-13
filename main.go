package main

import (
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

	var long = []int{20, 21, 22, 23, 24, 25}
	var short = []int{7, 8, 9, 10, 11, 12}

	// double boucle pour trouver meilleurs

	// lab.go pour analyser les données

	for i, _ := range long {

		testStrategy := Strategy{
			Asset:     "ETHUSDC",
			Intervals: []Interval{m5, m15, m30, h1, h2, h4},
			Filters:   []Filter{},
			Main: Signal{
				Name: "EMA",
				Type: "Moving Average",
				Params: map[Indicator]int{
					SMA_short: short[i],
					SMA_long:  long[i],
				},
			},
		}
		result := testStrategy.StrategyTester(client)
		result.AppendToHistory()

	}

}
