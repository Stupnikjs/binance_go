package main

var PAIRS = []string{
	"BTCUSDC", "ETHUSDC", "LINKUSDC", "ALGOUSDC", "HBARUSDC", "SOLUSDC", "AAVEUSDC"}

type ITTrader interface {
	Buy() error
	Sell() error
}
