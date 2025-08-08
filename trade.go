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

type Condition func(client *binance_connector.Client, pair string, interval string) bool

type Strategy struct {
	Asset         string
	Amount        float64
	BuyCondition  Condition
	SellCondition Condition
	*Trade
}

type Trade struct {
	Buy_price  *float64
	Buy_time   int64
	Sell_price *float64
	Sell_time  int64
	Status     int // 0 for not initiate // 1 for buy but not sell // 2 for sell
}

func (s *Strategy) Buy(client *binance_connector.Client) error {
	if s.Trade == nil {
		response, err := client.NewCreateOrderService().
			Symbol(s.Asset).
			Side("BUY").
			Type("MARKET").
			Quantity(s.Amount).
			Do(context.Background())

		// extract response to a struct
		var orderResponse CreateOrderResponse
		jsonBytes, err := json.Marshal(response)
		err = json.Unmarshal(jsonBytes, &orderResponse)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON to struct: %v", err)
		}
		trade := Trade{}
		s.Trade = &trade
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		s.Trade.Buy_price = &float_price
		s.Trade.Buy_time = orderResponse.TransactTime
		s.Trade.Status = 1

		return err

	}
	return nil
}
func (s *Strategy) Sell(client *binance_connector.Client) error {
	if s.Trade.Status == 1 {
		response, err := client.NewCreateOrderService().
			Symbol(s.Asset).
			Side("SELL").
			Type("MARKET").
			Quantity(s.Amount).
			Do(context.Background())

		// extract response to a struct
		var orderResponse CreateOrderResponse
		jsonBytes, err := json.Marshal(response)
		err = json.Unmarshal(jsonBytes, &orderResponse)
		if err != nil {
			log.Fatalf("Failed to unmarshal JSON to struct: %v", err)
		}
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		s.Trade.Sell_price = &float_price
		s.Trade.Sell_time = orderResponse.TransactTime
		return err

	}
	return nil
}

func (s *Strategy) GetGain(client *binance_connector.Client) (float64, error) {
	if s.Trade.Buy_price == nil || s.Trade.Sell_price == nil {
		return 0, fmt.Errorf("Trade not closed")
	}
	gain := *s.Trade.Sell_price*s.Amount - *s.Trade.Buy_price*s.Amount

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
