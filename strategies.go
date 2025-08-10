package main

import (
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)



type StrategyStat struct {
	Asset     string
	StratName string
	Interval  string
	Ratio     float64
}


// Strategy Tester 
/*


*/


func (stat *StrategyStat) StrategyTester(client *binance_connector.Client, smallPeriod int, bigPeriod int) {
	small_sma_field := fmt.Sprintf("sma_%d", smallPeriod)
	big_sma_field := fmt.Sprintf("sma_%d", bigPeriod)

	type trade struct {
		buyPrice  float64
		buyStamp  int
		sellPrice float64
		sellStamp int
	}
	closedTrade := []trade{}
	klines := IndicatorstoKlines(GetKlines(client, stat.Asset, stat.Interval, 1000), smallPeriod, bigPeriod, 14)

	var bigOverSmallPrev bool
	t := trade{}
	for _, k := range klines {
		if _, ok := k.Indicators[small_sma_field]; !ok {
			continue
		}
		if _, ok := k.Indicators[big_sma_field]; !ok {
			continue
		}

		bigOverSmall := k.Indicators[small_sma_field] < k.Indicators[big_sma_field]
		if !bigOverSmall && bigOverSmallPrev {
			f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
			if err != nil {
				fmt.Println(err)
			}
			t.buyPrice = f_close
			t.buyStamp = int(k.Kline_binance.CloseTime)

		}
		if bigOverSmall && !bigOverSmallPrev && t.buyStamp != 0 {
			f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
			if err != nil {
				fmt.Println(err)
			}
			t.sellPrice = f_close
			t.sellStamp = int(k.Kline_binance.CloseTime)
			closedTrade = append(closedTrade, t)
			t = trade{}
		}
		bigOverSmallPrev = k.Indicators[small_sma_field] < k.Indicators[big_sma_field]
	}

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.sellPrice / t.buyPrice) * prev_ratio
		prev_ratio = ratio

	}
	fmt.Println(ratio)
}

// underOver is positive if u check if RSI is above underOver
// its negative if u check if RSI is under
