package main

import (
	"log"
	"os"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

var testStrategy = Strategy{
	Amount:    0.01,
	Asset:     "ETHUSDC",
	Intervals: []Interval{m5, m15, m30, h1, h2, h4},
	Main: Signal{
		Name: "EMA",
		Type: "Moving Average",
		Params: map[Indicator]int{
			SMA_short: 9,
			SMA_long:  25,
		},
	},
}

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

	/*
		res := testStrategy.StrategyTester(client)
		if err != nil {
			log.Fatal(err)
		}
		res.AppendToHistory()
	*/
	FindBestMAParams(client, "XRPUSDC", []Interval{m5, m15, m30, h1, h2, h4})
	// AnalyseReport("XRPUSDC")

}
