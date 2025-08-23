package main

import (
	"context"
	"fmt"
	"strconv"
	"sync"

	"github.com/Stupnikjs/binance_go/order"
	binance_connector "github.com/binance/binance-connector-go"
)

type Strategy struct {
	USDCAmount float64
	Type       string
	Params     IndicatorsParams
	Intervals  []Interval
}

type Result struct {
	Pair       string
	StartStamp int
	EndStamp   int
	Ratio      float64
	Params     IndicatorsParams
}

var PARAMS = IndicatorsParams{

	short_period_MA: 9,
	long_period_MA:  21,
	super_long_MA:   200,
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

func OverSuperLong(kline *Klines, i int) bool {
	f_close, err := strconv.ParseFloat(kline.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f_close > kline.Indicators[SMA_super_long][i]
}

func (s *Strategy) TestWrapper(client *binance_connector.Client) ([]Result, error) {
	var wg sync.WaitGroup
	resultsChan := make(chan Result, len(PAIRS)) // Buffered channel
	var results []Result
	wg.Add(len(PAIRS))
	go func() {
		for result := range resultsChan {
			results = append(results, result)
		}
	}()
	for _, p := range PAIRS {
		go func(p string) {
			defer wg.Done()
			amount := ConvertUSDCtoPAIR(client, s.USDCAmount, p)
			klines := IndicatorstoKlines(
				client,
				p,
				s.Intervals,
				PARAMS)
			result := InitResult(p, klines[0].Array)
			result.Params = s.Params
			closedTrade := []BackTestTrader{}
			var prev bool
			t := InitBackTestTrader(p, amount, klines[0])
			loop := t.LoopBuilder(*s)
			for i := 0; i < len(klines[0].Indicators[SMA_super_long])-1; i++ {
				curr, err := loop(klines[0], &prev, i)
				if err != nil {
					fmt.Println(err)
				}
				prev = curr

				if t.TradeOver {
					closedTrade = append(closedTrade, *t)
					t = InitBackTestTrader(p, amount, klines[0])

				}

			}
			result.Ratio = 1.0
			for _, t := range closedTrade {
				result.Ratio = t.Sell_price / t.Buy_price
			}
			resultsChan <- result
		}(p)

	}
	go func() {
		wg.Wait()
		close(resultsChan)
	}()
	for result := range resultsChan {
		results = append(results, result)
	}
	return results, nil
}

func (s *Strategy) ParralelRunWrapper(ctx context.Context, client *binance_connector.Client) ([]*LiveTrader, error) {
	var wg sync.WaitGroup
	var mu sync.Mutex // Mutex to protect the shared slice

	// Create a channel to send completed trades to the main goroutine.
	tradesChan := make(chan *LiveTrader, len(PAIRS))

	wg.Add(len(PAIRS))

	for _, p := range PAIRS {
		// Launch a new goroutine for each currency pair.
		go func(p string) {
			defer wg.Done()

			// A true live trading loop would be more sophisticated. For this example,
			// we'll just run the logic once. A real application would need to
			// handle a continuous event stream (e.g., from a WebSocket).

			amount := ConvertUSDCtoPAIR(client, s.USDCAmount, p)
			order.PrintUSDCBalance(client)

			// Initialize a live trader for this specific pair
			t := InitLiveTrader(p, amount, client)
			prev := false
			loop := t.LoopBuilder(*s)
			_ = prev
			_ = loop
			for {
				select {
				case <-ctx.Done():
					// Context was canceled, so we exit the goroutine gracefully.
					fmt.Printf("Shutting down live trading for %s\n", p)
					return
				default:
					// Continue with the live trading logic
				}

				klines := IndicatorstoKlines(client, p, s.Intervals, PARAMS)

				// Check for empty klines data to prevent panic
				if len(klines) == 0 {
					fmt.Printf("Error: No klines data for pair %s\n", p)
					return
				}

				curr, err := loop(klines[0], &prev, len(klines[0].Array)-1)

				if err != nil {
					fmt.Printf("Error in loop for %s: %v\n", p, err)
					return
				}

				// If the trade is over, send it to the channel.
				if t.TradeOver {
					tradesChan <- t
				}
				prev = curr

			}

		}(p)
	}

	// This goroutine waits for all workers to finish and then closes the channel.
	go func() {
		wg.Wait()
		close(tradesChan)
	}()

	var completedTrades []*LiveTrader
	// Collect results from the channel until it's closed.
	for trade := range tradesChan {
		mu.Lock()
		completedTrades = append(completedTrades, trade)
		mu.Unlock()
	}

	return completedTrades, nil
}
