package main

var PAIRS = []string{
	"BTCUSDC", "ETHUSDC", "LINKUSDC",
}

type ITTrader interface {
	Buy() error
	Sell() error
	LoopBuilder(s Strategy) func(klines *Klines, prevOver *bool, i int) (bool, error) // change so it works with go routines
	// We'll also add a method to check if the trade is over.
	IsTradeOver() bool
	GetGain() (error, float64)
	SetStop() error
}
