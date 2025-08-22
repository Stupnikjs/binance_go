package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

type Wrapper struct {
	Asset     string
	Amount    float64
	Intervals []Interval // first interval is the main interval where we want to trade market
	Main      Signal
}

type Signal struct {
	Name   string
	Type   string
	Params map[Indicator]int
}

type Result struct {
	Pair       string
	StartStamp int
	EndStamp   int
	Ratio      float64
	Params     map[Indicator]int
}

func (s *Wrapper) InitResult(pair string, klines []*binance_connector.KlinesResponse) Result {
	result := Result{}
	result.StartStamp = int(klines[0].CloseTime)
	result.EndStamp = int(klines[len(klines)-1].CloseTime)
	result.Pair = pair
	return result
}

func (r *Result) GetRatioLive(trades []*LiveTrader) {
	r.Ratio = 1
	for _, t := range trades {
		r.Ratio = (t.Sell_price / t.Buy_price) * r.Ratio
	}
}

func (s *Wrapper) SetupParams() IndicatorsParams {
	return IndicatorsParams{
		short_period_MA: s.Main.Params[SMA_short],
		long_period_MA:  s.Main.Params[SMA_long],
		super_long_MA:   s.Main.Params[SMA_super_long],
		RSI_coef:        s.Main.Params[RSI],
	}
}

func OverSuperLong(kline *Klines, i int) bool {
	f_close, err := strconv.ParseFloat(kline.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f_close > kline.Indicators[SMA_super_long][i]
}

func (s *Wrapper) Test(client *binance_connector.Client) (*Result, error) {

	// setup klines
	params := s.SetupParams()
	klines := IndicatorstoKlines(
		client,
		s.Asset,
		s.Intervals,
		params)

	result := s.InitResult(s.Asset, klines[0].Array)
	closedTrade := []BackTestTrader{}
	var prev bool
	t := InitBackTestTrader(s.Asset, s.Amount, klines[0])
	strat := Strategy{}
	loop := t.LoopBuilder(strat)
	for i := 0; i < len(klines[0].Indicators[SMA_super_long])-1; i++ {
		curr, err := loop(klines[0], &prev, i)
		if err != nil {
			return nil, err
		}
		prev = curr

		if t.TradeOver {
			closedTrade = append(closedTrade, *t)
			t = InitBackTestTrader(s.Asset, s.Amount, klines[0])

		}

	}
	for _, t := range closedTrade {
		fmt.Println(t.Buy_price, t.Sell_price)
	}
	return &result, nil
}

func (s *Wrapper) RunSetup(client *binance_connector.Client) (IndicatorsParams, Result, *LiveTrader, Strategy) {
	fmt.Println("-- STARTING -- ")
	params := s.SetupParams()
	result := Result{}
	result.Ratio = 1
	t := InitLiveTrader(s.Asset, s.Amount, client)
	strat := Strategy{
		Type:   "Cross Over EMA",
		Params: params,
	}
	return params, result, t, strat
}

func (s *Wrapper) Run(client *binance_connector.Client) (*Result, error) {

	// setup
	tradeOver := []*LiveTrader{}
	PrintUSDCBalance(client)
	params, result, t, strat := s.RunSetup(client)
	prev := false
	loop := t.LoopBuilder(strat)
	for len(tradeOver) < 6 {
		klines := IndicatorstoKlines(client, s.Asset, s.Intervals, params)

		// curr is true is long period is over small
		curr, err := loop(klines[0], &prev, len(klines[0].Array)-1)
		if err != nil {
			return nil, err
		}

		prev = curr
		if t.TradeOver {
			tradeOver = append(tradeOver, t)
			PrintUSDCBalance(client)
			t = InitLiveTrader(s.Asset, s.Amount, client)
		}
		duration, err := IntervalToTime(s.Intervals[0])
		if err != nil {
			return nil, err
		}
		time.Sleep(duration)
	}
	result.GetRatioLive(tradeOver)
	return &result, nil
}

// decomposer fonction

func GetAssetBalance(client *binance_connector.Client, asset string) (float64, error) {

	account, err := client.NewGetAccountService().Do(context.Background())
	for i := range account.Balances {
		if asset == account.Balances[i].Asset {
			amount, err := strconv.ParseFloat(account.Balances[i].Free, 64)
			return amount, err
		}
	}
	return 0, err
}

func PrintUSDCBalance(client *binance_connector.Client) {
	usdc, err := GetAssetBalance(client, "USDC")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("USDC: %f \n", usdc)
}
