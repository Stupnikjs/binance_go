package main

import "github.com/google/uuid"

var PAIRS = []string{
	"BTCUSDC", "ETHUSDC", "LINKUSDC", "ALGOUSDC", "HBARUSDC", "SOLUSDC", "AAVEUSDC"}

type ITTrader interface {
	Buy() error
	Sell() error
}

// Exemple
type Trade struct {
	Id        uuid.UUID
	BuyPrice  float64
	BuyTime   int
	SellPrice float64
	SellTime  int
}
