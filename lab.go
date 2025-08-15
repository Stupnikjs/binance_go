package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"
	"sync"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

func FindBestMAParams(client *binance_connector.Client, pair string, interval []Interval) {
	var allResult []StrategyResult = make([]StrategyResult, 32*16)
	s := Strategy{
		Asset:     pair,
		Amount:    0.1,
		Intervals: interval,
		Main: Signal{
			Name:   "EMA",
			Type:   "Moving Average",
			Params: make(map[Indicator]int),
		},
	}
	m := Strategy{
		Asset:     pair,
		Amount:    0.1,
		Intervals: interval,
		Main: Signal{
			Name:   "SMA",
			Type:   "Moving Average",
			Params: make(map[Indicator]int),
		},
	}
	for i := 4; i < 20; i++ {
		for j := 8; j < 40; j++ {
			start_iter := time.Now().UnixMicro()
			if i >= j {
				continue
			} // i is < to j now
			s.Main.Params[SMA_short] = i
			s.Main.Params[SMA_long] = j
			m.Main.Params[SMA_short] = i
			m.Main.Params[SMA_long] = j
			r := s.StrategyTester(client)
			l := m.StrategyTester(client)
			allResult = append(allResult, r)
			allResult = append(allResult, l)
			end_iter := time.Now().UnixMicro()
			fmt.Printf("loop took %v \n", (end_iter - start_iter))
		}
	}
	filename := fmt.Sprintf("%s_report.json", strings.ToLower(pair))
	SaveJsonResult(filename, allResult)
}

func ParralelFindBestMAParams(client *binance_connector.Client, pair string, interval []Interval) {
	var allResult []StrategyResult
	var wg sync.WaitGroup
	resultsChan := make(chan StrategyResult, 32*16) // Buffered channel

	// Step 1: Launch a collector goroutine
	go func() {
		for result := range resultsChan {
			allResult = append(allResult, result)
		}
	}()

	for i := 4; i < 20; i++ {
		for j := 8; j < 40; j++ {
			if i >= j {
				continue
			}

			wg.Add(1) // Increment the WaitGroup counter
			go func(short, long int) {
				defer wg.Done() // Decrement the counter when the goroutine finishes

				// Create copies of the strategies for each goroutine
				s := Strategy{
					Asset:     pair,
					Amount:    0.1,
					Intervals: interval,
					Main: Signal{
						Name:   "EMA",
						Type:   "Moving Average",
						Params: make(map[Indicator]int),
					},
				}
				m := Strategy{
					Asset:     pair,
					Amount:    0.1,
					Intervals: interval,
					Main: Signal{
						Name:   "SMA",
						Type:   "Moving Average",
						Params: make(map[Indicator]int),
					},
				}

				s.Main.Params[SMA_short] = short
				s.Main.Params[SMA_long] = long
				m.Main.Params[SMA_short] = short
				m.Main.Params[SMA_long] = long

				r := s.StrategyTester(client)
				l := m.StrategyTester(client)

				resultsChan <- r // Send result to the channel
				resultsChan <- l // Send result to the channel
			}(i, j) // Pass i and j as arguments to the goroutine
		}
	}

	// Wait for all goroutines to finish
	wg.Wait()

	// Close the channel after all goroutines are done
	close(resultsChan)

	// Collect results from the channel
	for result := range resultsChan {
		allResult = append(allResult, result)
	}

	filename := fmt.Sprintf("%s_report.json", strings.ToLower(pair))
	SaveJsonResult(filename, allResult)
	AnalyseReport(pair)
}

func AnalyseReport(pairname string) {

	fileName := fmt.Sprintf("%s_report.json", strings.ToLower(pairname))
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Println(err)

	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		fmt.Println(err)

	}
	var results []StrategyResult
	err = json.Unmarshal(bytes, &results)
	if err != nil {
		fmt.Println(bytes)

	}
	slices.SortFunc(results, func(a, b StrategyResult) int {
		if a.Ratio > b.Ratio {
			return -1
		}
		if a.Ratio < b.Ratio {
			return 1
		}
		return 0
	})

	for _, r := range results[:20] {
		fmt.Println(r.Ratio, r.Strategy.Main)
	}

}
