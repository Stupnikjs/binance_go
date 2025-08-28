package klines

import (
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

type Indicator string

const (
	PRICE Indicator = "PRICE"
	RSI   Indicator = "RSI"
)

var Indicators = []Indicator{
	PRICE, RSI,
}

func IntervalToTime(interval Interval) (time.Duration, error) {
	// The time.ParseDuration function can handle strings like "10m", "1h", "1m30s".
	// It returns a duration and an error.
	return time.ParseDuration(string(interval))
}

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type IndicatorMapFunc map[Indicator]func([]float64, int) []float64

func GetRSI(klines []binance_connector.KlinesResponse, rsi_coef float64) []float64 {
	return analysis.RSIcalc(CloseFromKlines(klines), int(rsi_coef))
}

func BuildKlineArrData(pair string, interval []Interval) []binance_connector.KlinesResponse {

	path := path.Join("data", string(interval[0]), strings.ToLower(pair))
	kline, err := LoadKlinesFromFile(path)
	if err != nil {
		fmt.Println(err)
	}

	return kline
}

func CloseFromKlines(klines []binance_connector.KlinesResponse) []float64 {
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

func VolumeFromKlines(klines []binance_connector.KlinesResponse) []float64 {
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
