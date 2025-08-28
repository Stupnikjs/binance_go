package main

import (
	"sync"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	"github.com/google/uuid"
)

type BackTestTrader struct {
	Pair     string
	Amounts  []float64
	Data     map[kli.Indicator][]float64
	BuyCond  func()
	SellCond func()
}

type TradeReport struct {
	TradeId    uuid.UUID
	Asset      string
	Amount     float64
	TradeOver  bool
	Buy_price  float64
	Buy_time   int64
	Sell_price float64
	Sell_time  int64
}

// Buy simulates a buy order.
func (t *BackTestTrader) Buy() error {

	return nil
}
func (b *BackTestTrader) Loop(wg *sync.WaitGroup, tr chan TradeReport) error {
	t := TradeReport{}
	defer wg.Done()
	for i := 0; i < len(b.Data); i++ {
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

/*
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


*/
