package main

import (
	"context"
	"fmt"

	binance_connector "github.com/binance/binance-connector-go"
)

type Trade struct {
	Asset      string
	Buy_price  *float64
	Buy_time   any
	Sell_price *float64
	Sell_time  any
	Amount     float64
}

func (t *Trade) Buy(client *binance_connector.Client) error {
	if t.Buy_price == nil {
		response, err := client.NewCreateOrderService().
			Symbol(t.Asset).
			Side("BUY").
			Type("MARKET").
			Quantity(t.Amount).
			Do(context.Background())

		fmt.Println(response, err)

	}
	return nil
}
