package main

import (
	"fmt"

	"github.com/google/uuid"
)

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

func (t *Trade) Ratio() float64 {
	if t.SellTime == 0 || t.BuyTime == 0 {
		fmt.Println("cant take ratio of unclosed trade")
		return 0.0
	}
	return t.SellPrice / t.BuyPrice
}
