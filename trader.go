package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

type BackTestTrader struct {
	Id           int64
	Client       *binance_connector.Client
	Asset        string
	Amount       float64
	IndicatorMap map[Indicator]float64
	TradeOver    bool
	Buy_price    float64
	Buy_time     int64
	Sell_price   float64
	Sell_time    int64
}
type LiveTrader struct {
	Id           int64
	Client       *binance_connector.Client
	Asset        string
	Amount       float64
	IndicatorMap map[Indicator]float64
	TradeOver    bool
	Buy_price    float64
	Buy_time     int64
	Sell_price   float64
	Sell_time    int64
}

type ITTrader interface {
	Buy(*Klines, int) error
	Sell(*Klines, int) error
}

func (t *LiveTrader) Buy() error {
	if !t.TradeOver {
		response, err := t.BuildOrder(t.Client, "BUY")

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
		t.Id = orderResponse.OrderID
		fmt.Printf("trade opened %v", t)
		return err

	}
	return nil
}

func (t *LiveTrader) Sell() error {
	if !t.TradeOver {
		response, err := t.Client.NewCreateOrderService().
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

func (t *BackTestTrader) Buy(k *Klines, i int) {
	f_close, err := strconv.ParseFloat(k.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	t.Buy_price = f_close
	t.Buy_time = int64(k.Array[i].CloseTime)
}

func (t *BackTestTrader) Sell(k *Klines, i int) {
	f_close, err := strconv.ParseFloat(k.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	t.Sell_price = f_close
	t.Sell_time = int64(k.Array[i].CloseTime)
	t.TradeOver = true

}
