package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

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

func (t *LiveTrader) Buy() error {

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

		fmt.Printf("trade closed with %f gains \n", t.Sell_price*t.Amount-t.Amount*t.Buy_price)
		return err

	}
	return nil
}

func (t *LiveTrader) Loop(klines *Klines, prevOver *bool, i int) (bool, error) {
	// Your CrossOver logic, adapted for the LiveTrader struct.
	closeOverMAsuperLong := OverSuperLong(klines, i)

	bigOverSmall := klines.Indicators[SMA_short][i] < klines.Indicators[SMA_long][i]

	if !bigOverSmall && *prevOver && closeOverMAsuperLong {
		if err := t.Buy(); err != nil {
			return false, err
		}
	}
	if bigOverSmall && !*prevOver && t.Buy_time != 0 {
		if err := t.Sell(); err != nil {
			return false, err
		}
	}
	return bigOverSmall, nil
}

// InitLiveTrader initializes a new LiveTrader instance.
func InitLiveTrader(pair string, amount float64, client *binance_connector.Client) *LiveTrader {
	return &LiveTrader{
		Asset:        pair,
		Amount:       amount,
		Client:       client,
		IndicatorMap: make(map[Indicator]float64),
	}
}
