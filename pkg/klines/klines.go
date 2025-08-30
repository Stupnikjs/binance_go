package klines

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	binance_connector "github.com/binance/binance-connector-go"
)

var Indicators = []Indicator{
	{"RSI", "Price", "Coef", 14, nil}, {
		"EMA", "Price", "Avg", 9, nil},
}

type Interval string

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type IndicatorMapFunc map[Indicator]func([]float64, int) []float64

func InitMapIndic(ind []Indicator) mapIndicator {
	var mapIndic mapIndicator = mapIndicator{
		"RSI": analysis.RSIcalc,
		"EMA": analysis.EMAcalc,
		"SMA": analysis.SMAcalc,
	}
	return mapIndic

}

const (
	m1  Interval = "1m"
	m5  Interval = "5m"
	m15 Interval = "15m"
	m30 Interval = "30m"
	h1  Interval = "1h"
	h2  Interval = "2h"
	h4  Interval = "4h"
)

type Indicator struct {
	Name   string
	Data   string
	Type   string
	Params int
	Values *[]float64
}

type mapIndicator map[string]func([]float64, int) []float64

type Data struct {
	Pair       string
	Interval   Interval
	Indicators []Indicator
	Volume     []float64
	Price      []float64
}

func (d *Data) InitData(k []*binance_connector.KlinesResponse, pair string, interval Interval, ind []Indicator) {
	mapInd := InitMapIndic(ind)
	d.Interval = interval
	d.Pair = pair
	d.Price = CloseFromKlines(k)
	d.Volume = VolumeFromKlines(k)
	for _, i := range ind {
		if i.Data == "Price" {
			*i.Values = mapInd[i.Name](d.Price, i.Params)
		}
		if i.Data == "Volume" {
			*i.Values = mapInd[i.Name](d.Volume, i.Params)
		}

	}

}

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
