package main

import (
	"fmt"
	"log"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

const (
	RSI Indicator = "RSI"
	SMA Indicator = "SMA"
	EMA Indicator = "EMA"
	VOL Indicator = "VOL"
)

type Indicator string

type Strategy struct {
	Asset    string
	Interval string
	Filters  []Filter
	Main     Signal
}

type Signal struct {
	Name   string
	Type   string
	Params map[string]int
}

type Filter struct {
	BuyFilter bool
	Indicator
	Value    float64
	Interval string
}

type StrategyResult struct {
	StartStamp int
	EndStamp   int
	Ratio      float64
	Strategy   *Strategy
}

func (strat *Strategy) StrategyTester(client *binance_connector.Client) StrategyResult {

	type trade struct {
		buyPrice  float64
		buyStamp  int
		sellPrice float64
		sellStamp int
	}

	closedTrade := []trade{}
	klinesClean := GetKlines(client, strat.Asset, strat.Interval, 1000)
	result := StrategyResult{}
	result.StartStamp = int(klinesClean[0].CloseTime)
	result.EndStamp = int(klinesClean[len(klinesClean)-1].CloseTime)
	if strat.Main.Type == "Moving Average" {
		klines := IndicatorstoKlines(
			klinesClean,
			strat.Main.Params["short"], strat.Main.Params["long"],
			14)
		if strat.Main.Name != "SMA" && strat.Main.Name != "EMA" {
			log.Fatal("wrong strat name ")
		}
		var small_field, big_field string
		if strat.Main.Name == "SMA" {
			small_field = fmt.Sprintf("sma_%d", strat.Main.Params["short"])
			big_field = fmt.Sprintf("sma_%d", strat.Main.Params["long"])
		}
		if strat.Main.Name == "EMA" {
			small_field = fmt.Sprintf("ema_%d", strat.Main.Params["short"])
			big_field = fmt.Sprintf("ema_%d", strat.Main.Params["long"])
		}

		var bigOverSmallPrev bool
		t := trade{}
		for _, k := range klines {
			if _, ok := k.Indicators[small_field]; !ok {
				continue
			}
			if _, ok := k.Indicators[big_field]; !ok {
				continue
			}

			bigOverSmall := k.Indicators[small_field] < k.Indicators[big_field]
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
			bigOverSmallPrev = k.Indicators[small_field] < k.Indicators[big_field]
		}

	}

	// modify to pass.Indicators

	// type == "Mooving Average"

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.sellPrice / t.buyPrice) * prev_ratio
		prev_ratio = ratio

	}
	result.Ratio = ratio
	result.Strategy = strat
	return result
}

// underOver is positive if u check if RSI is above underOver
// its negative if u check if RSI is under

// systeme de mail 

func (strat *Strategy) StrategyApply(client *binance_connector.Client) {
    
   

   // create loop 

   
   for {




   } 
   // create limit 
   // build klines
   // check for buy signal 
   // build trade 
   // check for sell signal 

}