package main

import (
	"fmt"
	"log"
	"os"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	"github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

var indic = []klines.Indicator{
	{Name: "RSI", Interval: klines.Interv[1], Type: "Price", Calculator: analysis.RSIcalc, Param: 14},
}

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
	k, err := GetSomeTestKlines()
	if err != nil {
		fmt.Println(err)
	}
	featured := klines.BuildFeaturedKlinesArray(k, indic)
	klines.FeaturedKlinesToCSV("test.csv", featured)
}

func GetSomeTestKlines() ([]*binance_connector.KlinesResponse, error) {
	return klines.LoadKlinesFromFile(klines.GetFilePathName("ALGOUSDC", klines.Interv[1]))
}

func SaveLastKlines(client *binance_connector.Client, intervals []klines.Interval) error {
	for _, i := range PAIRS {
		err := klines.AppendNewData(client, i, intervals)
		if err != nil {
			return err
		}

	}
	return nil
}
