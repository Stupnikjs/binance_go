package main

import (
	"strconv"

	"github.com/Stupnikjs/binance_go/pkg/klines"
)

type BackTestTrader struct {
	Pair       string
	Curr       *Trade
	Indicators []klines.Indicator
}

// Buy simulates a buy order.
func (t *BackTestTrader) Buy() error {

	return nil
}
func (t *BackTestTrader) Sell() error {

	return nil
}

func (b *BackTestTrader) SetStop(price float64) error {
	return nil
}

type BackTestResult struct {
	Pair  string
	Ratio float64
}

func InitBackTestTrader(pair string, indicators []klines.Indicator) BackTestTrader {
	return BackTestTrader{
		Pair:       pair,
		Indicators: indicators,
		Curr:       nil,
	}
}

func InitBackTestResult() BackTestResult {
	return BackTestResult{}
}

func BackTestTradesToResult(trades []Trade) BackTestResult {
	ratio := 1.0
	for _, t := range trades {
		ratio = (t.SellPrice - t.BuyPrice) / t.BuyPrice * ratio
	}
	return InitBackTestResult()
}

func (b *BackTestTrader) Iterate(feature klines.FeaturedKlines, prev *bool) *Trade {
	// for EMA cross over
	shortOverLong := feature.FeaturesMap[b.Indicators[0].GetMapKey()] > feature.FeaturesMap[b.Indicators[1].GetMapKey()]
	if shortOverLong && *prev {
		// Buy Logic
		b.Buy()
		f_price, err := strconv.ParseFloat(feature.Close, 64)
		if err != nil {
			panic(err)
		}
		tradeNew := InitTrade(f_price, int(feature.CloseTime))
		b.Curr = &tradeNew
	}
	if !shortOverLong && !*prev {
		// Sell Logic
		b.Sell()
		return b.Curr
	}
	return nil
}
