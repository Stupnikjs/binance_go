package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

func (t *Trader) Buy(client *binance_connector.Client) error {
	if !t.TradeOver {
		response, err := t.BuildOrder(client, "BUY")

		if err != nil {
			return err
		}
		orderResponse, err := t.ParseResponse(response)
		if err != nil {
			return err
		}

		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		if err != nil {
			return err
		}
		t.Buy_price = float_price
		t.Buy_time = orderResponse.TransactTime
		fmt.Printf("trade opened %v", t)
		return err

	}
	return nil
}

func (t *Trader) Sell(client *binance_connector.Client) error {
	if !t.TradeOver {
		response, err := client.NewCreateOrderService().
			Symbol(t.Asset).
			Side("SELL").
			Type("MARKET").
			Quantity(t.Amount).
			Do(context.Background())
		if err != nil {
			return err
		}
		// extract response to a struct
		var orderResponse CreateOrderResponse
		jsonBytes, err := json.Marshal(response)
		if err != nil {
			return err
		}
		err = json.Unmarshal(jsonBytes, &orderResponse)
		if err != nil {
			return err
		}
		if len(orderResponse.Fills) == 0 {
			return fmt.Errorf(" error buying asset ")
		}
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		if err != nil {
			return err
		}

		t.Sell_price = float_price
		t.Sell_time = orderResponse.TransactTime

		return err

	}
	return nil
}
