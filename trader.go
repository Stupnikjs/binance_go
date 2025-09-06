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

func InitTrade(buyPrice float64, buyTime int) Trade {
	return Trade{
		Id:       uuid.New(),
		BuyPrice: buyPrice,
		BuyTime:  buyTime,
	}
}
