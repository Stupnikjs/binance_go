package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"

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

type Trader struct {
	Client          *binance_connector.Client
	Asset           string
	Amount          float64
	IndicatorMap    map[Indicator]float64
	TradeInProgress bool
	Buy_price       float64
	Buy_time        int64
	Sell_price      float64
	Sell_time       int64
}

type StrategyResult struct {
	Pair       string
	StartStamp int
	EndStamp   int
	Ratio      float64
	Params     map[Indicator]int
}

func InitTrader(pair string, amount float64) Trader {
	return Trader{
		Asset:           pair,
		Amount:          amount,
		TradeInProgress: false,
		IndicatorMap:    make(map[Indicator]float64),
	}
}
func (s *Strategy) InitResult(pair string, klines []*binance_connector.KlinesResponse) StrategyResult {
	result := StrategyResult{}
	result.StartStamp = int(klines[0].CloseTime)
	result.EndStamp = int(klines[len(klines)-1].CloseTime)
	result.Pair = pair
	return result
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

func (s *Strategy) StrategyTester(client *binance_connector.Client) StrategyResult {
	params := IndicatorsParams{
		short_period_MA: s.Main.Params[SMA_short],
		long_period_MA:  s.Main.Params[SMA_long],
		super_long_MA:   s.Main.Params[SMA_super_long],
		RSI_coef:        s.Main.Params[RSI],
	}
	klines := IndicatorstoKlines(
		client,
		s.Asset,
		s.Intervals,
		params)

	result := s.InitResult(s.Asset, klines[0].Array)
	closedTrade := []Trader{}

	if s.Main.Type == "Moving Average" {
		if s.Main.Name != "SMA" && s.Main.Name != "EMA" {
			log.Fatal("wrong strat name ")
		}

		var bigOverSmallPrev bool
		t := InitTrader(s.Asset, s.Amount)
		offset := s.Main.Params[SMA_super_long]
		for i := 1; i < len(klines[0].Indicators[SMA_super_long]); i++ {
			// check crossOver or Under
			// implement checking or RSI in the timeframe in klines[1] or klines[2]
			err := MeltRSIKline(klines[0], klines[2])
			if err != nil {
				fmt.Println(err)
			}
			// pass index offset
			bigOverSmallPrev = t.CrossMA(klines[0], i, offset, &bigOverSmallPrev, &closedTrade, true)
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

// underOver is positive if u check if RSI is above underOver
// its negative if u check if RSI is under

// systeme de mail

func MeltRSIKline(receiver *Klines, origin *Klines) error {
	RSI_origin := origin.Indicators[RSI]
	var targetIndicator Indicator
	// switch on Intervals
	switch origin.Interval {
	case m15:
		targetIndicator = RSI_15m
	case m30:
		targetIndicator = RSI_30m
	case h1:
		targetIndicator = RSI_1h
	case h2:
		targetIndicator = RSI_2h
	case h4:
		targetIndicator = RSI_4h
	default:
		return fmt.Errorf("no indicator valid found")
	}
	n := 0
	for i := range origin.Array {
		curr := origin.Array[n]
		if receiver.Array[i].CloseTime < curr.CloseTime {
			receiver.Indicators[targetIndicator][i] = RSI_origin[i]
		}

	}
	return nil
}

func (s *Strategy) StrategyApply(client *binance_connector.Client) error {
	params := IndicatorsParams{
		short_period_MA: s.Main.Params[SMA_short],
		long_period_MA:  s.Main.Params[SMA_long],
		super_long_MA:   s.Main.Params[SMA_super_long],
		RSI_coef:        s.Main.Params[RSI],
	}
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

		// fetch new kline
		klines := IndicatorstoKlines(client, s.Asset, s.Intervals, params)
		// err := MeltRSIKline(klines[0], klines[3]) NOT USED YET
		bigOverSmall := t.CrossMA(klines[0], len(klines[0].Array)-1, 0, &bigOverSmallPrev, &tradeOver, false)

		if !bigOverSmall && bigOverSmallPrev && !t.TradeInProgress {

			err = t.Buy(client)
			fmt.Printf("Buying at %v %v \n", t.Buy_time, t.Buy_price)
			if err != nil {
				return err
			}
		}
		if bigOverSmall && !bigOverSmallPrev && t.TradeInProgress {
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
			t = Trader{
				Asset:        s.Asset,
				Amount:       s.Amount,
				IndicatorMap: map[Indicator]float64{},
			}

		}
		bigOverSmallPrev = bigOverSmall

	}
	filename := fmt.Sprintf("%s_trade_report.json", s.Asset)
	SaveJsonTrader(filename, tradeOver)
	return nil
}

// decomposer fonction

func (t *Trader) BuildOrder(client *binance_connector.Client) (interface{}, error) {
	return client.NewCreateOrderService().
		Symbol(t.Asset).
		Side("BUY").
		Type("MARKET").
		Quantity(t.Amount).
		Do(context.Background())

}

func (t *Trader) ParseResponse(response interface{}) (*CreateOrderResponse, error) {
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
func (t *Trader) Buy(client *binance_connector.Client) error {
	if !t.TradeInProgress {
		response, err := t.BuildOrder(client)

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
		t.TradeInProgress = true
		fmt.Printf("trade opened %v", t)
		return err

	}
	return nil
}

func (t *Trader) BuyTest(k *Klines, offset int) {
	f_close, err := strconv.ParseFloat(k.Array[offset].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	t.Buy_price = f_close
	t.Buy_time = int64(k.Array[offset].CloseTime)
}

func (t *Trader) Sell(client *binance_connector.Client) error {
	if t.TradeInProgress {
		response, err := client.NewCreateOrderService().
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
		t.TradeInProgress = false

		return err

	}
	return nil
}

func (t *Trader) SellTest(k *Klines, offset int, closed *[]Trader) {
	f_close, err := strconv.ParseFloat(k.Array[offset].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	t.Sell_price = f_close
	t.Sell_time = int64(k.Array[offset].CloseTime)
	*closed = append(*closed, *t)
	*t = InitTrader(t.Asset, t.Amount)
}

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
