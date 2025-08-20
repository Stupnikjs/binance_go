package main

import (
	"fmt"
	"log"
	"os"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

var testStrat Strategy = Strategy{
	Asset:     "ETHUSDC",
	Amount:    0.002,
	Intervals: []Interval{m1, m5, m15, m30, h1},
	Main: Signal{
		Name:   "EMA",
		Type:   "Moving Average",
		Params: make(map[Indicator]int),
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

	testStrat.Main.Params[SMA_short] = 9
	testStrat.Main.Params[SMA_long] = 15
	testStrat.Main.Params[SMA_super_long] = 200

	t, err := testStrat.Run(client)
	if err != nil {
		fmt.Println(err)
	}

	_ = t

	// Get API credentials from environment variables
	/*
		err = FetchReports(client, []Interval{m5, m15, m30, h1})
		if err != nil {
			fmt.Println(err)
		}

		GiveReportData(m5)
	*/
}

// creer des commandes pour backtester
