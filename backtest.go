package main

import (
	"fmt"
	"strconv"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/google/uuid"
)

type BackTestTrader struct {
	Id           uuid.UUID
	Client       *binance_connector.Client
	Asset        string
	Amount       float64
	IndicatorMap map[kli.Indicator]float64
	TradeOver    bool
	Buy_price    float64
	Buy_time     int64
	Sell_price   float64
	Sell_time    int64
	Klines       *kli.Klines // Needed for backtesting
	Index        int         // Current index in the Klines array
}

func InitBackTestTrader(pair string, amount float64, klines *kli.Klines) *BackTestTrader {
	return &BackTestTrader{
		Id:           uuid.New(),
		Asset:        pair,
		Amount:       amount,
		IndicatorMap: make(map[kli.Indicator]float64),
		Klines:       klines,
		Index:        0,
		TradeOver:    false,
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

func (b *BackTestTrader) LoopBuilder(s Strategy) func(klines *kli.Klines, prevOver *bool, i int) (bool, error) {
	return func(klines *kli.Klines, prevOver *bool, i int) (bool, error) {
		b.Klines = klines
		b.Index += 1

		closeOverMAsuperLong := OverSuperLong(klines, i)
		bigOverSmall := klines.Indicators[kli.EMA_short][i] < klines.Indicators[kli.EMA_long][i]

		if !bigOverSmall && *prevOver && closeOverMAsuperLong && b.Buy_time == 0 {
			if err := b.Buy(); err != nil {
				return false, err
			}
		}
		if bigOverSmall && !*prevOver && b.Buy_time != 0 && b.Sell_time == 0 {
			if err := b.Sell(); err != nil {
				return false, err
			}

		}
		// its not sell signal
		return bigOverSmall, nil
	}
}

func (b *BackTestTrader) SetStop(price float64) error {
	return nil
}
func (b *BackTestTrader) GetGain() (float64, error) {
	if !b.TradeOver {
		return 0, fmt.Errorf("trade still in progress %v ", b)
	}
	return b.Sell_price*b.Amount - b.Buy_price*b.Amount, nil
}
