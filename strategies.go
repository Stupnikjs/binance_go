package main

import (
	"fmt"
	"strconv"
	"sync"

	kli "github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
)

type Strategy struct {
	USDCAmount float64
	Type       string
	Params     kli.IndicatorsParams
	Intervals  []kli.Interval
}

type Result struct {
	Pair       string
	StartStamp int
	EndStamp   int
	Ratio      float64
	Params     kli.IndicatorsParams
}

var PARAMS = kli.IndicatorsParams{
	Short_period_MA: 9,
	Long_period_MA:  21,
	Super_long_MA:   200,
	RSI_coef:        14,
	VROC_coef:       15,
}

func InitResult(pair string, klines []*binance_connector.KlinesResponse) Result {
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

func OverSuperLong(kline *kli.Klines, i int) bool {
	f_close, err := strconv.ParseFloat(kline.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f_close > kline.Indicators[kli.SMA_super_long][i]
}

func (s *Strategy) RunWrapper(client *binance_connector.Client) ([]LiveTrader, error) {
	var wg sync.WaitGroup
	closedTradeChan := make(chan LiveTrader, len(PAIRS)) // Buffered channel
	closedTrade := []LiveTrader{}
	wg.Add(len(PAIRS) * 20)
	go func() {
		for trade := range closedTradeChan {
			closedTrade = append(closedTrade, trade)
		}
	}()
	for _, p := range PAIRS {
		// init a trader and send it for each go routines
		amout := kli.ConvertUSDCtoPAIR(client, s.USDCAmount, p)
		t := InitLiveTrader(p, amout, client)
		// THERE
		go t.RoutineWrapper(&wg, closedTradeChan, p, s.Intervals)
	}

	wg.Wait()
	close(closedTradeChan)

	return closedTrade, nil
}
