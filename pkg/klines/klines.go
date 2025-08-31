package klines

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

func BuildKlineArrData(pair string, interval []Interval) []*binance_connector.KlinesResponse {

	path := path.Join("data", string(interval[0]), strings.ToLower(pair))
	kline, err := LoadKlinesFromFile(path)
	if err != nil {
		fmt.Println(err)
	}

	return kline
}

func FetchKlines(client *binance_connector.Client, pair string, intervals []Interval) ([]*binance_connector.KlinesResponse, error) {
	return client.NewKlinesService().
		Symbol(pair).
		Interval(string(intervals[0])).
		Limit(1000).
		Do(context.Background())

}

func CloseFromKlines(klines []*binance_connector.KlinesResponse) []float64 {
	closingPrices := make([]float64, len(klines))
	for i, kline := range klines {
		f_close, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			fmt.Println(err)
		}
		closingPrices[i] = f_close

	}
	return closingPrices
}

func VolumeFromKlines(klines []*binance_connector.KlinesResponse) []float64 {
	volumes := make([]float64, len(klines))
	for i, kline := range klines {
		f_vol, err := strconv.ParseFloat(kline.Volume, 64)
		if err != nil {
			fmt.Println(err)
		}
		volumes[i] = f_vol

	}
	return volumes
}
