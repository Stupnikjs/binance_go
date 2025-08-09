package main

import (
	"fmt"

	binance_connector "github.com/binance/binance-connector-go"
)

func RSIbuyCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, -25)
}
func RSIsellCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, 75)
}

func EMAcrossOver(client *binance_connector.Client, pair string, EMAcoefShort int, EMAcoefLong int, over bool) bool {
	klines := GetKlines(client, pair, "1m", EMAcoefLong*2)
	closing := CloseFromKlines(klines)
	EMAshort := ExponetialMovingAverage(closing, EMAcoefShort)
	EMALong := ExponetialMovingAverage(closing, EMAcoefLong)

	bullishCross := EMAshort[len(EMAshort)-2] < EMALong[len(EMALong)-2] && EMAshort[len(EMAshort)-1] > EMALong[len(EMALong)-1]
	bearishCross := EMAshort[len(EMAshort)-2] > EMALong[len(EMALong)-2] && EMAshort[len(EMAshort)-1] < EMALong[len(EMALong)-1]
	if over {
		if bullishCross {
			return true
		}
	}
	if !over {
		if bearishCross {
			return true
		}
	}
	return false

}

type StrategyStat struct {
	Asset     string
	StratName string
	Interval  string
	Ratio     float64
}

func (stat *StrategyStat) SMATest(client *binance_connector.Client) {

	type trade struct {
		buyPrice  float64
		sellPrice float64
	}
	closedTrade := []trade{}
	klines := GetKlines(client, stat.Asset, stat.Interval, 1000)
	close := CloseFromKlines(klines)

	sma_9 := SMA(close, 9)
	sma_20 := SMA(close, 20)

	sma_9_alligned := sma_9[11:]
	close_alligned := close[19:]

	prev := sma_9_alligned[0] >= sma_20[0]

	curTrade := trade{}
	for i := 1; i < len(close_alligned); i++ {

		crossover := sma_9_alligned[i] > sma_20[i]
		crossunder := sma_9_alligned[i] < sma_20[i]
		if crossover && !prev && curTrade.buyPrice == 0 {
			curTrade.buyPrice = close_alligned[i]

		}
		if crossunder && prev && curTrade.buyPrice != 0 {
			curTrade.sellPrice = close_alligned[i]
			closedTrade = append(closedTrade, curTrade)
			curTrade = trade{}

		}
		prev = crossover
	}
	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.sellPrice / t.buyPrice) * prev_ratio
		prev_ratio = ratio

	}
	fmt.Println(ratio)
}
