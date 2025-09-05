package main

import (
	"sync"

	"github.com/Stupnikjs/binance_go/pkg/klines"
)

type BackTestTrader struct {
	Pair     string
	Curr     *Trade
}

// Buy simulates a buy order.
func (t *BackTestTrader) Buy() error {

	return nil
}
func (b *BackTestTrader) Loop(wg *sync.WaitGroup, tr chan Trade) error {
	defer wg.Done()
	for i := 0; i < len(b.Amounts); i++ {
		// logic buy => fill tradeReport

		// initTrade
		// logic sell

		// if trade over

		tr <- t

	}
	return nil
}

func (b *BackTestTrader) SetStop(price float64) error {
	return nil
}

type BackTestResult {
	Pair string 
	Ratio float64
	
}

func BackTestTradesToResult(trades []Trade) BackTestResult() {
	ratio := 1.0
	for _ , t := range trades {
		ratio = (t.SellPrice - t.BuyPrice) / t.BuyPrice * ratio  
		}
	return initBackTestResult(pair, ratio) 
}



func (b *BackTestTrader) Iterate(feature klines.FeaturedKlines, prev *bool) *Trade {
	// for EMA cross over 
	shortOverLong := feature.FeaturesMap[b.Indicators[0].GetKey()] > feature.FeaturesMap[b.Indicators[1].GetKey()]
	if shortOverLong && *prev {
		// Buy Logic
		b.Buy()
		b.Curr = intitTrade()
		  }
	if !shortOverLong && !*prev  {
		// Sell Logic 
		b.Sell()
		return b.Curr
		}
	
	
}


func RunBackTest() error {
	var wg sync.WaitGroup
	reports := []TradeReport{}
	tradeReports := make(chan TradeReport, 1000)
	wg.Add(len(PAIRS))
	go func() {
		for _, r := range tradeReports {
			reports = append(reports, r)
		}
	}()
	for _, p := range PAIRS {
		// init BACKTESTTRADER
		b := InitBackTestTrader()
		go b.Loop(wg, tradeReport) {

		}()
	}
	wg.Wait()

}



