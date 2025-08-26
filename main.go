package main

import (
	"fmt"
	"log"
	"os"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var apiKey = os.Getenv("API_KEY")
	var secretKey = os.Getenv("SECRET_KEY")

	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")
	_ = client

	if err != nil {
		fmt.Println(err)
	}
	/*
		testStrat := Strategy{
			USDCAmount: 10,
			Type:       "Cross over EMA",
			Params: kli.IndicatorsParams{
				Short_period_MA: 9,
				Long_period_MA:  21,
				RSI_coef:        14,
				VROC_coef:       15,
			},
			Intervals: []kli.Interval{m5, m15, h1, h2},
		}

		_, err = testStrat.RunWrapper(client)
		if err != nil {
			fmt.Println(err)
		}

	*/

	lines, err := kli.LoadKlinesFromFile("BTCUSDC", kli.Interv[1])

	for _, i := range lines {
		fmt.Println(i)
	}

	if err != nil {
		fmt.Println(err)
	}

	// Get API credentials from environment variable

}

// creer des commandes pour backte
