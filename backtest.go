package main

import (
	"sync"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	"github.com/google/uuid"
)

type BackTestTrader struct {
	Pair     string
	Amounts  []float64
	Curr     *Trade
	Data     map[kli.Indicator][]float64
	BuyCond  func()
	SellCond func()
}

type Trade struct {
	Id        uuid.UUID
	Asset     string
	Amount    float64
	TradeOver bool
}

type TradeReport struct {
	TradeId    uuid.UUID
	Asset      string
	Amount     float64
	Buy_price  float64
	Buy_time   int64
	Sell_price float64
	Sell_time  int64
}

// Buy simulates a buy order.
func (t *BackTestTrader) Buy() error {

	return nil
}
func (t *BackTestTrader) Loop(chan TradeReport) error {

	return nil
}

func (b *BackTestTrader) SetStop(price float64) error {
	return nil
}

func RunBackTest() error {
	var wg sync.WaitGroup
	reports := []TradeReport{}
	tradeReports := make(chan TradeReport, 1000)
	wg.Add(50)
	go func() {
		for _, r := range tradeReports {
			reports = append(reports, r)
		}
	}()
	for _, p := range PAIRS {
		// init BACKTESTTRADER
		b := InitBackTestTrader()
		go func(wg, tradeReport) {

		}()
	}
	wg.Wait()

}
