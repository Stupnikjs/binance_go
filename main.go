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
	_ = klines.BuildFeaturedKlinesArray(k, indic)
	err = SaveLastKlines(client, klines.Interv[1:])
	fmt.Println(err)
}

func PairLoop(pair string, ind []klines.Indicator) {
	k := klines.LoadKlinesFromFile(klines.FileName(pair, klines.Interv[1:]))
	featured := klines.BuildFeaturedKlinesArray(k, ind)
	prev := false

	trades := []Trade{}
	b := InitBackTestTrader(pair, ind)
	// ind[0] is short [1] is big
	for i, f := range featured {
		t := b.Iterate(f, &prev)
		if t != nil {
			trades = append(trades, *t)
		}

	}

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
