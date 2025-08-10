package main

// Price unity is USDC for now

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

type Condition func(client *binance_connector.Client, pair string, interval string) bool

type Trader struct {
	Asset           string
	Amount          float64
	BuyCondition    Condition `json:"-"`
	SellCondition   Condition `json:"-"`
	TradeInProgress bool
	*Trade
}

type Trade struct {
	Buy_price  *float64
	Buy_time   int64
	Sell_price *float64
	Sell_time  int64
}

// decomposer fonction

func (s *Trader) BuildOrder(client *binance_connector.Client) (interface{}, error) {
	return client.NewCreateOrderService().
		Symbol(s.Asset).
		Side("BUY").
		Type("MARKET").
		Quantity(s.Amount).
		Do(context.Background())

}

func (s *Trader) ParseResponse(response interface{}) (*CreateOrderResponse, error) {
	var orderResponse CreateOrderResponse
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &orderResponse)
	if err != nil {
		return nil, err
	}
	if len(orderResponse.Fills) == 0 {
		return nil, fmt.Errorf(" received empty response ")
	}
	return &orderResponse, nil
}
func (s *Trader) Buy(client *binance_connector.Client) error {
	if !s.TradeInProgress {
		response, err := s.BuildOrder(client)

		if err != nil {
			return err
		}
		orderResponse, err := s.ParseResponse(response)
		if err != nil {
			return err
		}

		trade := Trade{}
		s.Trade = &trade
		float_price, err := strconv.ParseFloat(orderResponse.Fills[0].Price, 64)
		if err != nil {
			return err
		}
		s.Trade.Buy_price = &float_price
		s.Trade.Buy_time = orderResponse.TransactTime
		s.TradeInProgress = true

		return err

	}
	return nil
}
func (s *Trader) Sell(client *binance_connector.Client) error {
	if s.TradeInProgress {
		response, err := client.NewCreateOrderService().
			Symbol(s.Asset).
			Side("SELL").
			Type("MARKET").
			Quantity(s.Amount).
			Do(context.Background())

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

		s.Trade.Sell_price = &float_price
		s.Trade.Sell_time = orderResponse.TransactTime
		s.TradeInProgress = false
		err = SaveTrade(*s)
		return err

	}
	return nil
}

func (s *Trader) GetGain(client *binance_connector.Client) (float64, error) {
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
