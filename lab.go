package main

import (
	"fmt"
	"strconv"
	"sync"

	binance_connector "github.com/binance/binance-connector-go"
)

var PAIRS = []string{
	"ADAUSDC",
	"ALGOUSDC",
	"BTCUSDC",
	"BNBUSDC",
	"XLMUSDC",
	"SUIUSDC",
	"ETHUSDC",
	"HBARUSDC",
	"LINKUSDC",
	"XRPUSDC",
}

/*
* list all reports
* and test best params
* then append to report
 */

func ConvertUSDCtoPAIR(client *binance_connector.Client, USDCamount float64, pair string) float64 {
	klines := BuildKlinesArr(client, pair, []Interval{m1})
	f_close, err := strconv.ParseFloat(klines[0].Array[0].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return USDCamount / f_close
}

/*
* create strategy with some for loop generated parameters
* then test it
 */
func ParralelTest(client *binance_connector.Client, interval []Interval) []Result {
	var allResult []Result
	var wg sync.WaitGroup
	resultsChan := make(chan Result, 32*16) // Buffered channel

	// Step 1: Launch a collector goroutine

	go func() {
		for result := range resultsChan {
			allResult = append(allResult, result)
		}
	}()
	for _, pair := range PAIRS {
		wg.Add(1) // Increment the WaitGroup counter
		go func(pair string) {
			defer wg.Done() // Decrement the counter when the goroutine finishes

			// Create copies of the strategies for each goroutine
			amount := ConvertUSDCtoPAIR(client, 30, pair)
			paramMap := make(map[Indicator]int)
			paramMap[SMA_short] = 13
			paramMap[SMA_long] = 43
			paramMap[SMA_super_long] = 200
			s := Wrapper{
				Asset:     pair,
				Amount:    amount,
				Intervals: interval,
				Main: Signal{
					Name:   "EMA",
					Type:   "Moving Average",
					Params: paramMap,
				},
			}

			s.Main.Params[SMA_short] = 13
			s.Main.Params[SMA_long] = 43
			s.Main.Params[SMA_super_long] = 200

			r, _ := s.Test(client)

			resultsChan <- *r // Send result to the channel

		}(pair) // Pass i and j as arguments to the goroutine

	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the channel after all goroutines are done
	close(resultsChan)

	// Collect results from the channel
	for result := range resultsChan {
		allResult = append(allResult, result)
	}

	return allResult

}
