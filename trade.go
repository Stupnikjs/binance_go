package main

// Price unity is USDC for now

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

		// extract response to a struct
		_ = response

		return err

	}
	return nil
}
func (t *Trade) Sell(client *binance_connector.Client) error {
	if t.Sell_price == nil {
		response, err := client.NewCreateOrderService().
			Symbol(t.Asset).
			Side("SELL").
			Type("MARKET").
			Quantity(t.Amount).
			Do(context.Background())

		// extract response to a struct
		_ = response

		return err

	}
	return nil
}

func (t *Trade) GetGain(client *binance_connector.Client) (float64, error) {
	if t.Buy_price == nil || t.Sell_price == nil {
		return 0, fmt.Errorf("Trade not closed")
	}
	gain := *t.Sell_price*t.Amount - *t.Buy_price*t.Amount

	return gain, nil
}
