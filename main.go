package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	"github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

var indic = []klines.Indicator{
	{Name: "EMA_short", Interval: klines.Interv[1], Type: "Price", Calculator: analysis.EMAcalc, Param: 9},
	{Name: "EMA_long", Interval: klines.Interv[1], Type: "Price", Calculator: analysis.EMAcalc, Param: 15},
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

	for {
		err = SaveLastKlines(client, klines.Interv[1:])
		if err != nil {
			fmt.Println(err)
		}
		time.Sleep(12 * time.Hour)
	}

}

func TestPairLoop() {
	tradeChan := make(chan Trade, 100)
	go PairLoop("BTCUSDC", indic, tradeChan)
	trades := []Trade{}

	defer func() {
		for t := range tradeChan {
			trades = append(trades, t)

		}
		result := BackTestTradesToResult(trades)
		fmt.Println(result)
	}()
}

func PairLoop(pair string, ind []klines.Indicator, tradeChan chan Trade) {

	k, _ := klines.LoadKlinesFromFile(klines.FileName(pair, klines.Interv))
	featured := klines.BuildFeaturedKlinesArray(k, ind)
	prev := false

	b := InitBackTestTrader(pair, ind)
	for _, f := range featured {
		t := b.Iterate(f, &prev)
		if t != nil {
			fmt.Println(t)
			tradeChan <- *t
		}

	}
	close(tradeChan)

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
		err = klines.CheckWholeHasNoTimeGap(i, intervals[0])
		if err != nil {
			return err
		}
	}

	return nil
}
