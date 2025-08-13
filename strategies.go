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
		small_field, long_field, err := s.GetMAFields()

		if err != nil {
			fmt.Println(err)
		}
		var bigOverSmallPrev bool
		t := s.InitTrader()
		for _, k := range klines {
			if _, ok := k.Indicators[small_field]; !ok {
				continue
			}
			if _, ok := k.Indicators[long_field]; !ok {
				continue
			}

			bigOverSmall := k.Indicators[small_field] < k.Indicators[long_field]
			if !bigOverSmall && bigOverSmallPrev {
				f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.Buy_price = f_close
				t.Buy_time = int64(k.Kline_binance.CloseTime)

			}
			if bigOverSmall && !bigOverSmallPrev && t.Buy_time != 0 {
				f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.Sell_price = f_close
				t.Sell_time = int64(k.Kline_binance.CloseTime)
				closedTrade = append(closedTrade, t)
				t = Trader{}
			}
			bigOverSmallPrev = k.Indicators[small_field] < k.Indicators[long_field]
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

	tradeOver := []Trader{}
	result := StrategyResult{}
	result.Ratio = 1

	t := Trader{
		Asset:           s.Asset,
		Amount:          s.Amount,
		TradeInProgress: false,
	}

	short_field, long_field, err := s.GetMAFields()
	if err != nil {
		return err
	}

	for result.Ratio > 0.5 {

		klineNative := GetKlines(client, s.Asset, s.Interval, 100) // limite must be > big period
		kline := IndicatorstoKlines(klineNative, s.Main.Params["short"], s.Main.Params["long"], 14)
		bearishPrev := kline[len(kline)-2].Indicators[short_field] < kline[len(kline)-2].Indicators[long_field]
		bullish := kline[len(kline)-2].Indicators[short_field] >= kline[len(kline)-2].Indicators[long_field]

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
