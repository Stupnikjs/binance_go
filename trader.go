package main

import kli "github.com/Stupnikjs/binance_go/pkg/klines"

var PAIRS = []string{
	"BTCUSDC", "ETHUSDC", "LINKUSDC", "ALGOUSDC", "HBARUSDC", "SOLUSDC", "AAVEUSDC"}

type ITTrader interface {
	Buy() error
	Sell() error
	LoopBuilder(s Strategy) func(klines *kli.Klines, prevOver *bool, i int) (bool, error) // change so it works with go routines
	// We'll also add a method to check if the trade is over.
	IsTradeOver() bool
	GetGain() (error, float64)
	SetStop() error
}
