package main

import (
	"fmt"
	"log"
	"os"
	"sync"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	"github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

var indic = []klines.Indicator{
	{Name: "EMA_short", Interval: klines.Interv[2], Type: "Price", Calculator: analysis.EMAcalc, Param: 5},
	{Name: "EMA_long", Interval: klines.Interv[2], Type: "Price", Calculator: analysis.EMAcalc, Param: 15},
	{Name: "RSI", Interval: klines.Interv[2], Type: "Price", Calculator: analysis.EMAcalc, Param: 14},
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

	bigArr := klines.GetEMAIndicatorsArray()

	max := 1.0
	for _, arr := range bigArr {
		ratio := GetAvgRatio(arr)
		if ratio > max {
			fmt.Println(arr, ratio)
			max = ratio
		}

	}

}

func GetAvgRatio(ind []klines.Indicator) (result float64) {

	tradeChan := make(chan Trade, 200)
	var wg sync.WaitGroup
	wg.Add(len(PAIRS))

	for _, p := range PAIRS {
		go EMASTRAT(p, indic, &wg, tradeChan)
	}
	rr := []Trade{}

	wg.Wait()
	close(tradeChan)

	for t := range tradeChan {

		rr = append(rr, t)
	}
	if len(rr) < 1 {
		return 1
	} else {
		sum := 0.0
		for _, r := range rr {
			sum += r.Ratio()
		}
		result = sum / float64(len(rr))
	}

	return result

}

func EMASTRAT(pair string, ind []klines.Indicator, wg *sync.WaitGroup, tradeChan chan Trade) {
	defer wg.Done()
	k, _ := klines.LoadKlinesFromFile(klines.FileName(pair, klines.Interv[2:]))
	featured := klines.BuildFeaturedKlinesArray(k, ind)
	prev := false
	currTrade := Trade{}
	for _, f := range featured {
		shortOverLong, err := f.EMAShortOverLong(ind)
		if err != nil {
			log.Fatal(err)
		}
		// Buy
		if shortOverLong && !prev && currTrade.BuyTime == 0 {
			fmt.Println(f.FeaturesMap[ind[2].GetMapKey()])
			currTrade = Trade{
				Id:       uuid.New(),
				BuyPrice: f.FloatClose(),
				BuyTime:  int(f.CloseTime),
			}
		}
		if !shortOverLong && prev && currTrade.BuyTime != 0 {
			currTrade.SellPrice = f.FloatClose()
			currTrade.SellTime = int(f.CloseTime)
			tradeChan <- currTrade
			currTrade = Trade{}
		}
		prev = shortOverLong
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
		err = klines.CheckWholeHasNoTimeGap(i, intervals[0])
		if err != nil {
			return err
		}
		fmt.Println("no time gap")
	}

	return nil
}
