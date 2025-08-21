package main

type ITTrader interface {
	Buy() error
	Sell() error
	Loop(klines *Klines, prevOver *bool, i int) (bool, error)
	// We'll also add a method to check if the trade is over.
	IsTradeOver() bool
	SetStop() error
}
