package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

type Indicator string
type Strategy struct {
	Asset     string
	Amount    float64
	Intervals []Interval // first interval is the main interval where we want to trade market
	Filters   []Filter
	Main      Signal
}

type Signal struct {
	Name   string
	Type   string
	Params map[Indicator]int
}

type Filter struct {
	BuyFilter bool
	Indicator
	Value    float64
	Interval string
}
type Trader struct {
	Asset           string
	Amount          float64
	TradeInProgress bool
	Buy_price       float64
	Buy_time        int64
	Sell_price      float64
	Sell_time       int64
}

func (s *Strategy) InitTrader() Trader {
	return Trader{
		Asset:           s.Asset,
		Amount:          s.Amount,
		TradeInProgress: false,
	}
}
func (s *Strategy) InitResult(klines []*binance_connector.KlinesResponse) StrategyResult {
	result := StrategyResult{}
	result.StartStamp = int(klines[0].CloseTime)
	result.EndStamp = int(klines[len(klines)-1].CloseTime)
	return result
}

func (s *Strategy) GetMAFields() (string, string, error) {

	short_field := fmt.Sprintf("%s_%d", strings.ToLower(s.Main.Name), s.Main.Params["short"])
	long_field := fmt.Sprintf("%s_%d", strings.ToLower(s.Main.Name), s.Main.Params["long"])
	var err error
	if short_field == "" || long_field == "" {
		err = fmt.Errorf("Error parsing MA field")
	}

	return short_field, long_field, err
}

type StrategyResult struct {
	StartStamp int // BORDERS OF THE SAMPLE TESTED
	EndStamp   int
	Ratio      float64
	Strategy   *Strategy
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
		RSI_coef:        s.Main.Params[RSI],
	}
	klines := IndicatorstoKlines(
		client,
		s.Asset,
		s.Intervals,
		params)

	result := s.InitResult(klines[0].Array)
	closedTrade := []Trader{}

	if s.Main.Type == "Moving Average" {
		if s.Main.Name != "SMA" && s.Main.Name != "EMA" {
			log.Fatal("wrong strat name ")
		}

		var bigOverSmallPrev bool
		t := s.InitTrader()
		index_long := s.Main.Params[SMA_long]
		for i := 0; i < len(klines[0].Array); i++ {
			if i < s.Main.Params[SMA_long] {
				continue
			}
			bigOverSmall := klines[0].Indicators[SMA_short][i] < klines[0].Indicators[SMA_long][i]
			if !bigOverSmall && bigOverSmallPrev {
				f_close, err := strconv.ParseFloat(klines[0].Array[i-index_long].Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.Buy_price = f_close
				t.Buy_time = int64(klines[0].Array[i-index_long].CloseTime)

			}
			if bigOverSmall && !bigOverSmallPrev && t.Buy_time != 0 {
				f_close, err := strconv.ParseFloat(klines[0].Array[i-index_long].Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.Sell_price = f_close
				t.Sell_time = int64(klines[0].Array[i-index_long].CloseTime)
				closedTrade = append(closedTrade, t)
				t = Trader{}
			}
			bigOverSmallPrev = bigOverSmall
		}

	}

	// modify to pass.Indicators

	// type == "Mooving Average"

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.Sell_price / t.Buy_price) * prev_ratio
		prev_ratio = ratio

	}
	result.Ratio = ratio
	result.Strategy = s
	return result
}

// underOver is positive if u check if RSI is above underOver
// its negative if u check if RSI is under

// systeme de mail

func (s *Strategy) StrategyApply(client *binance_connector.Client) error {
	params := IndicatorsParams{
		short_period_MA: s.Main.Params[SMA_short],
		long_period_MA:  s.Main.Params[SMA_long],
		RSI_coef:        s.Main.Params[RSI],
	}
	tradeOver := []Trader{}
	result := StrategyResult{}
	result.Ratio = 1

	t := Trader{
		Asset:           s.Asset,
		Amount:          s.Amount,
		TradeInProgress: false,
	}

	for result.Ratio > 0.5 {
		klines := IndicatorstoKlines(client, s.Asset, s.Intervals, params)
		SMA_short := klines[0].Indicators[SMA_short]
		SMA_long := klines[0].Indicators[SMA_long]
		bearishPrev := SMA_short[len(SMA_short)-2] < SMA_long[len(SMA_long)-2]
		bullish := SMA_short[len(SMA_short)-1] < SMA_long[len(SMA_long)-1]

		if bearishPrev && bullish && t.Buy_time == 0 {
			err := t.Buy(client)
			if err != nil {

				return err
			}
		}
		if !bearishPrev && !bullish && t.Buy_time != 0 {
			err := t.Sell(client)
			if err != nil {
				return err
			}
			tradeOver = append(tradeOver, t)
			t = Trader{}
		}

	}
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

		return err

	}
	return nil
}
func (t *Trader) Sell(client *binance_connector.Client) error {
	if t.TradeInProgress {
		response, err := client.NewCreateOrderService().
			Symbol(t.Asset).
			Side("SELL").
			Type("MARKET").
			Quantity(t.Amount).
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

		t.Sell_price = float_price
		t.Sell_time = orderResponse.TransactTime
		t.TradeInProgress = false

		return err

	}
	return nil
}

func (t *Trader) GetGain(client *binance_connector.Client) (float64, error) {
	if t.Buy_price == 0 || t.Sell_price == 0 {
		return 0, fmt.Errorf("Trade not closed")
	}
	gain := t.Sell_price*t.Amount - t.Buy_price*t.Amount

	return gain, nil
}

func StrategyLab() {}

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
