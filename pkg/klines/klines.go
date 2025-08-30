package klines

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	binance_connector "github.com/binance/binance-connector-go"
)

type Interval string

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
}

type mapIndicator map[string]func([]float64, int) []float64

var Indicators = []Indicator{

	{"RSI", "Price", "Coef", 14}, {
		"EMA", "Price", "Avg", 9},
}

func InitMapIndic(ind []Indicator) mapIndicator {
	var mapIndic mapIndicator = mapIndicator{
		"RSI": analysis.RSIcalc,
		"EMA": analysis.EMAcalc,
		"SMA": analysis.SMAcalc,
	}
	return mapIndic

}

func BuildSuperArray(k []*binance_connector.KlinesResponse, indicators []Indicator, sliceLen int) [][]float64 {
	var superArray [][]float64
	close := CloseFromKlines(k)
	vols := VolumeFromKlines(k)
	index := len(close) - sliceLen + 1
	superArray = append(superArray, close[index:])
	superArray = append(superArray, vols[index:])
	mapIndicator := InitMapIndic(indicators)
	for i := 2; i < len(indicators)-1; i++ {
		indicator := indicators[i-2]
		var indicArr []float64
		if indicator.Data == "Price" {
			indicArr = mapIndicator[indicator.Name](close, indicator.Params)
		}
		if indicator.Data == "Volume" {
			indicArr = mapIndicator[indicator.Name](vols, indicator.Params)
		}

		sliceIndex := len(indicArr) - sliceLen + 1
		superArray[i] = indicArr[sliceIndex:]
	}
	return superArray
}

func IntervalToTime(interval Interval) (time.Duration, error) {
	// The time.ParseDuration function can handle strings like "10m", "1h", "1m30s".
	// It returns a duration and an error.
	return time.ParseDuration(string(interval))
}

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type IndicatorMapFunc map[Indicator]func([]float64, int) []float64



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
