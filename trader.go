package main

type ITTrader interface {
	Buy() error
	Sell() error
	LoopBuilder(s Strategy) func(klines *Klines, prevOver *bool, i int) (bool, error)
	// We'll also add a method to check if the trade is over.
	IsTradeOver() bool
	GetGain() (error, float64)
	SetStop() error
}
type Strategy struct {
	Type   string
	Params IndicatorsParams
}
