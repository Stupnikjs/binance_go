package main

import (
	"context"
	"log"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

func GetKlines(client *binance_connector.Client, pair string, interval string, limit int) []*binance_connector.KlinesResponse {
	klines, err := client.NewKlinesService().
		Symbol(pair).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	return klines
}

func CloseFromKlines(klines []*binance_connector.KlinesResponse) []float64 {
	closingPrices := make([]float64, len(klines))
	for i, kline := range klines {
		f_close, err := strconv.ParseFloat(kline.Close, 64)
		closingPrices[i] = f_close
		if err != nil {
			log.Fatal(err)
		}
	}
	return closingPrices
}

func SMA(closingPrices []float64, period int) []float64 {
	var SMA []float64
	closingPriceSlice := closingPrices
	if len(closingPrices) < period {
		return SMA
	}
	for i := period - 1; i < len(closingPriceSlice); i++ {
		var sma float64
		slice := closingPriceSlice[i-period+1 : i+1]
		for _, n := range slice {
			sma += n

		}
		SMA = append(SMA, sma/float64(period))

	}
	return SMA
}
