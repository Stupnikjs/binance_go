package main

// Price unity is USDC for now

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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
		var orderResponse CreateOrderResponse
		jsonBytes, err := json.Marshal(response)
		err = json.Unmarshal(jsonBytes, &orderResponse)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON to struct: %v", err)
		}
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		t.Buy_price = &float_price
		t.Buy_time = orderResponse.TransactTime

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
		var orderResponse CreateOrderResponse
		jsonBytes, err := json.Marshal(response)
		err = json.Unmarshal(jsonBytes, &orderResponse)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON to struct: %v", err)
		}
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		t.Sell_price = &float_price
		t.Sell_time = orderResponse.TransactTime
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

type CreateOrderResponse struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Fills         []Fill `json:"fills"`
}

// Fill represents a single fill of the order.
type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}
