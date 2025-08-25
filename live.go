package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	"github.com/Stupnikjs/binance_go/pkg/order"
	binance_connector "github.com/binance/binance-connector-go"
)

type LiveTrader struct {
	Id           int64
	Client       *binance_connector.Client
	Asset        string
	Amount       float64
	IndicatorMap map[kli.Indicator]float64
	TradeOver    bool
	Buy_price    float64
	Buy_time     int64
	Sell_price   float64
	Sell_time    int64
	StopPrice    float64
}

func (t *LiveTrader) Buy() error {
	response, err := order.BuildOrder(t.Client, "BUY", t.Asset, t.Amount)
	if err != nil {
		return err
	}
	orderResponse, err := order.ParseResponse(response)
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
		var orderResponse order.CreateOrderResponse
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

func (t *LiveTrader) LoopBuilder() func(klines *kli.Klines, prevOver *bool, i int) (bool, error) {
	return func(klines *kli.Klines, prevOver *bool, i int) (bool, error) {
		if len(klines.Indicators[kli.EMA_short]) < 700 || len(klines.Indicators[kli.EMA_long]) < 700 {
			return false, fmt.Errorf("klines malformed")
		}
		bigOverSmall := klines.EMAShortOverLong(i)

		fmt.Printf("time: %s sma short : %f sma long: %f VROC: %f \n",
			order.TimeStampToDateString(int(klines.Array[i].CloseTime))[:20], klines.Indicators[kli.EMA_short][i],
			klines.Indicators[kli.EMA_long][i], klines.Indicators[kli.VROC][i])

		f_close, err := strconv.ParseFloat(klines.Array[i-1].Close, 64)
		if err != nil {
			return false, err
		}
		if t.StopPrice >= f_close {
			t.TradeOver = true
			fmt.Printf("trade closed %s %f \n", order.TimeStampToDateString(int(t.Sell_time)), t.Sell_price)
			return bigOverSmall, nil
		}
		if !bigOverSmall && *prevOver {

			if err := t.Buy(); err != nil {
				fmt.Printf("trade open %s %f \n", order.TimeStampToDateString(int(t.Buy_time)), t.Buy_price)
				err = t.SetStop(f_close)
				return false, err
			}

		}
		if bigOverSmall && !*prevOver && t.Buy_time != 0 {
			if err := t.Sell(); err != nil {
				return true, err
			}
		}
		return bigOverSmall, nil
	}
}

func (t *LiveTrader) SetStop(price float64) error {
	t.StopPrice = price
	return order.BuildStopLoss(t.Client, price, t.Asset, t.Amount)

}

// InitLiveTrader initializes a new LiveTrader instance.
func InitLiveTrader(pair string, amount float64, client *binance_connector.Client) *LiveTrader {
	return &LiveTrader{
		Asset:        pair,
		Amount:       amount,
		Client:       client,
		IndicatorMap: make(map[kli.Indicator]float64),
	}
}

func (t *LiveTrader) RoutineWrapper(wg *sync.WaitGroup, closedTradeChan chan LiveTrader, p string, Intervals []kli.Interval) error {
	defer wg.Done()
	prev := false
	loop := t.LoopBuilder()
	for !t.TradeOver {
		klines := kli.IndicatorstoKlines(t.Client, t.Asset, Intervals, PARAMS)

		signal, err := loop(klines[0], &prev, len(klines[0].Array)-1)
		if err != nil {
			return err
		}
		prev = signal
		timeToSleep, err := kli.IntervalToTime(Intervals[0])
		time.Sleep(timeToSleep)
		if err != nil {
			return err
		}
	}
	closedTradeChan <- *t
	return nil
}
