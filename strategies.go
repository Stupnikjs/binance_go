package main

import (
	"context"
	"fmt"
	"strconv"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

type Strategy struct {
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

type StrategyResult struct {
	Pair       string
	StartStamp int
	EndStamp   int
	Ratio      float64
	Params     map[Indicator]int
}

func (s *Strategy) InitResult(pair string, klines []*binance_connector.KlinesResponse) StrategyResult {
	result := StrategyResult{}
	result.StartStamp = int(klines[0].CloseTime)
	result.EndStamp = int(klines[len(klines)-1].CloseTime)
	result.Pair = pair
	return result
}

func (s *Strategy) SetupParams() IndicatorsParams {
	return IndicatorsParams{
		short_period_MA: s.Main.Params[SMA_short],
		long_period_MA:  s.Main.Params[SMA_long],
		super_long_MA:   s.Main.Params[SMA_super_long],
		RSI_coef:        s.Main.Params[RSI],
	}
}

func (s *Strategy) Test(client *binance_connector.Client) StrategyResult {
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

	for i := 0; i < len(klines[0].Indicators[SMA_super_long])-1; i++ {
		curr, _ := t.Loop(klines[0], &prev, i)
		prev = curr
		if t.TradeOver {
			closedTrade = append(closedTrade, *t)
			t = InitBackTestTrader(s.Asset, s.Amount, klines[0])

		}

	}

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.Sell_price / t.Buy_price) * prev_ratio
		prev_ratio = ratio

	}
	fmt.Println(len(closedTrade))
	result.Ratio = ratio
	result.Params = s.Main.Params
	return result
}

func (s *Strategy) Run(client *binance_connector.Client) ([]LiveTrader, error) {

	// setup
	params := s.SetupParams()
	tradeOver := []LiveTrader{}
	result := StrategyResult{}
	result.Ratio = 1
	t := InitLiveTrader(s.Asset, s.Amount, client)
	prev := false
	for len(tradeOver) < 10 {
		klines := IndicatorstoKlines(client, s.Asset, s.Intervals, params)
		curr, _ := t.Loop(klines[0], &prev, len(klines[0].Array)-1)
		prev = curr
		if t.TradeOver {
			tradeOver = append(tradeOver, *t)
			t = InitLiveTrader(s.Asset, s.Amount, client)
		}
		time.Sleep(IntervalToTime(s.Intervals[0]))
	}

	return tradeOver, nil
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
