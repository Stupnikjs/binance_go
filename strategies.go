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

const (
	RSI Indicator = "RSI"
	SMA Indicator = "SMA"
	EMA Indicator = "EMA"
	VOL Indicator = "VOL"
)

type Indicator string

type Strategy struct {
	Asset    string
	Amount   float64
	Interval string
	Filters  []Filter
	Main     Signal
}

type Signal struct {
	Name   string
	Type   string
	Params map[string]int
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

func (strat *Strategy) StrategyTester(client *binance_connector.Client) StrategyResult {

	type trade struct {
		buyPrice  float64
		buyStamp  int
		sellPrice float64
		sellStamp int
	}

	closedTrade := []trade{}
	klinesClean := GetKlines(client, strat.Asset, strat.Interval, 1000)
	result := StrategyResult{}
	result.StartStamp = int(klinesClean[0].CloseTime)
	result.EndStamp = int(klinesClean[len(klinesClean)-1].CloseTime)
	if strat.Main.Type == "Moving Average" {
		klines := IndicatorstoKlines(
			klinesClean,
			strat.Main.Params["short"], strat.Main.Params["long"],
			14)
		if strat.Main.Name != "SMA" && strat.Main.Name != "EMA" {
			log.Fatal("wrong strat name ")
		}
		var small_field, big_field string
		if strat.Main.Name == "SMA" {
			small_field = fmt.Sprintf("sma_%d", strat.Main.Params["short"])
			big_field = fmt.Sprintf("sma_%d", strat.Main.Params["long"])
		}
		if strat.Main.Name == "EMA" {
			small_field = fmt.Sprintf("ema_%d", strat.Main.Params["short"])
			big_field = fmt.Sprintf("ema_%d", strat.Main.Params["long"])
		}

		var bigOverSmallPrev bool
		t := trade{}
		for _, k := range klines {
			if _, ok := k.Indicators[small_field]; !ok {
				continue
			}
			if _, ok := k.Indicators[big_field]; !ok {
				continue
			}

			bigOverSmall := k.Indicators[small_field] < k.Indicators[big_field]
			if !bigOverSmall && bigOverSmallPrev {
				f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.buyPrice = f_close
				t.buyStamp = int(k.Kline_binance.CloseTime)

			}
			if bigOverSmall && !bigOverSmallPrev && t.buyStamp != 0 {
				f_close, err := strconv.ParseFloat(k.Kline_binance.Close, 64)
				if err != nil {
					fmt.Println(err)
				}
				t.sellPrice = f_close
				t.sellStamp = int(k.Kline_binance.CloseTime)
				closedTrade = append(closedTrade, t)
				t = trade{}
			}
			bigOverSmallPrev = k.Indicators[small_field] < k.Indicators[big_field]
		}

	}

	// modify to pass.Indicators

	// type == "Mooving Average"

	prev_ratio := 1.0
	ratio := 1.0
	for _, t := range closedTrade {
		ratio = (t.sellPrice / t.buyPrice) * prev_ratio
		prev_ratio = ratio

	}
	result.Ratio = ratio
	result.Strategy = strat
	return result
}

// underOver is positive if u check if RSI is above underOver
// its negative if u check if RSI is under

// systeme de mail

func (s *Strategy) StrategyApply(client *binance_connector.Client) {

	// create loop
	result := StrategyResult{}
	result.Ratio = 1

	short_field := fmt.Sprintf("%s_%d", strings.ToLower(s.Main.Name), s.Main.Params["short"])
	long_field := fmt.Sprintf("%s_%d", strings.ToLower(s.Main.Name), s.Main.Params["long"])

	for result.Ratio > 0.5 {

		klineNative := GetKlines(client, s.Asset, s.Interval, 100) // limite must be > big period
		kline := IndicatorstoKlines(klineNative, s.Main.Params["short"], s.Main.Params["long"], 14)
		bearishPrev := kline[len(kline)-2].Indicators[short_field] < kline[len(kline)-2].Indicators[long_field]
		bullish := kline[len(kline)-2].Indicators[short_field] >= kline[len(kline)-2].Indicators[long_field]

		t := Trader{
			Asset:           s.Asset,
			Amount:          s.Amount,
			TradeInProgress: false,
		}
		if bearishPrev && bullish && t.Buy_time == 0 {
			err := t.Buy(client)
			if err != nil {
				fmt.Println(err)
			}
		}
		if !bearishPrev && !bullish && t.Buy_time != 0 {
			err := t.Sell(client)
			if err != nil {
				fmt.Println(err)
			}
			// SAVE TRADE
			// reset trade
		}

	}
	// create limit
	// build klines
	// check for buy signal
	// build trade
	// check for sell signal

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
