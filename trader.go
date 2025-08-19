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
	Klines       *Klines // Needed for backtesting
	Index        int     // Current index in the Klines array
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
	Buy() error
	Sell() error
	Loop(klines *Klines, prevOver *bool, i int) (bool, error)
	// We'll also add a method to check if the trade is over.
	IsTradeOver() bool
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

// Buy simulates a buy order.
func (t *BackTestTrader) Buy() error {
	f_close, err := strconv.ParseFloat(t.Klines.Array[t.Index].Close, 64)
	if err != nil {
		return err
	}
	t.Buy_price = f_close
	t.Buy_time = int64(t.Klines.Array[t.Index].CloseTime)
	fmt.Printf("Backtest Buy at %.2f\n", t.Buy_price)
	return nil
}
func (t *BackTestTrader) Sell() error {
	f_close, err := strconv.ParseFloat(t.Klines.Array[t.Index].Close, 64)
	if err != nil {
		return err
	}
	t.Sell_price = f_close
	t.Sell_time = int64(t.Klines.Array[t.Index].CloseTime)
	t.TradeOver = true
	fmt.Printf("Backtest Sell at %.2f\n", t.Sell_price)
	return nil
}

func (t *BackTestTrader) Loop(klines *Klines, prevOver *bool, i int) (bool, error) {
	t.Klines = klines
	t.Index = i

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

func InitBackTestTrader(pair string, amount float64, klines *Klines) *BackTestTrader {
	return &BackTestTrader{
		Asset:        pair,
		Amount:       amount,
		IndicatorMap: make(map[Indicator]float64),
		Klines:       klines,
	}
}
