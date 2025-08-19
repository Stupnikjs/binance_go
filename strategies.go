package main

import (
	"context"
	"fmt"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

type Strategy struct {
	Asset     string
	Amount    float64
	Looper    func(*Trader, *Klines, *bool, int) bool
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

func InitBackTestTrader() BackTestTrader {
	return BackTestTrader{}
}
func InitLiveTrader() LiveTrader {
	return LiveTrader{}
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
	var bigOverSmallPrev bool
	t := InitTrader(s.Asset, s.Amount)
	fmt.Println(len(klines[0].Array), len(klines[0].Indicators[SMA_long]), len(klines[0].Indicators[SMA_short]))
	for i := 0; i < len(klines[0].Indicators[SMA_super_long]); i++ {
		_ = MeltRSIKline(klines[0], klines[2])
		bigOverSmall := s.Looper(&t, klines[0], &bigOverSmallPrev, i)
		bigOverSmallPrev = bigOverSmall
		if t.TradeOver {
			closedTrade = append(closedTrade, t)
			t = InitTrader(s.Asset, s.Amount)
		}

	}

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.Sell_price / t.Buy_price) * prev_ratio
		prev_ratio = ratio

	}
	result.Ratio = ratio
	result.Params = s.Main.Params
	return result
}

func (s *Strategy) Run(client *binance_connector.Client) error {

	// setup
	params := s.SetupParams()
	tradeOver := []Trader{}
	result := StrategyResult{}
	result.Ratio = 1
	t := InitTrader(s.Asset, s.Amount)
	var prevRatio = 1.0
	oldBalance, err := GetAssetBalance(client, "USDC")
	if err != nil {
		return err
	}

	var bigOverSmallPrev bool
	for result.Ratio < 1.2 && result.Ratio > 0.8 {
		// fetch new kline each loop turn
		klines := IndicatorstoKlines(client, s.Asset, s.Intervals, params)
		i := len(klines[0].Array) - 1

		// compare last items of SMA
		overSuperLong := OverSuperLong(klines[0], i)
		bigOverSmall := klines[0].Indicators[SMA_short][i] <
			klines[0].Indicators[SMA_long][i]

		if !bigOverSmall && bigOverSmallPrev && t.Buy_time == 0 && overSuperLong {
			err = t.BuyFuncs.Func(t, client)
			// placer un stop loss
			fmt.Printf("Buying at %v %v \n", t.Buy_time, t.Buy_price)
			if err != nil {
				return err
			}
		}
		if bigOverSmall && !bigOverSmallPrev && t.Buy_time != 0 {
			err := t.Sell(client)
			fmt.Printf("Selling at %v %v \n", t.Sell_time, t.Sell_price)
			if err != nil {
				return err
			}
			newBalance, err := GetAssetBalance(client, "USDC")
			if err != nil {
				return err
			}
			fmt.Printf("balance USDC: %v \n", newBalance-oldBalance)

			tradeOver = append(tradeOver, t)
			ratio := (t.Sell_price - t.Buy_price) * prevRatio
			result.Ratio = ratio
			fmt.Printf("Ratio : %v \n", ratio)
			prevRatio = ratio
			t = InitTrader(s.Asset, s.Amount)

		}
		bigOverSmallPrev = bigOverSmall

	}
	filename := fmt.Sprintf("%s_trade_report.json", s.Asset)
	SaveJsonTrader(filename, tradeOver)
	return nil
}

func CrossOver(t *ITTrader, klines *Klines, prev *bool, i int) bool {
	closeOverMAsuperLong := OverSuperLong(klines, i)
	bigOverSmall := klines.Indicators[SMA_short][i] < klines.Indicators[SMA_long][i]
	if !bigOverSmall && *prev && closeOverMAsuperLong {
		t.Buy(klines, i)
	}
	if bigOverSmall && !*prev && t.Buy_time != 0 {
		t.SellTest(klines, i)
	}
	return bigOverSmall
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
