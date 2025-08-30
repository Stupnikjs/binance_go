package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Stupnikjs/binance_go/pkg/klines"
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

	// Get API credentials from environment variable

	for _, i := range PAIRS {

		err = klines.AppendNewData(client, i, klines.Interv[1:])

		if err != nil {
			fmt.Println(err)
		}
	}

}



func BuildSignalArray([]klines.Indicator, )


func Loop(client *binance_connector.Client, pair string, interval []klines.Interval) {
	signal := "crossOver"
	for {
		klines, _  := klines.FetchKlines(client, pair, interval)
		// build two last items 
		signalArrays := buildSignalArray()
		if  

	}
	





}