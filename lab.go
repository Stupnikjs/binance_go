package main

import (
	"fmt"
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

func ParralelFindBestMAParams(client *binance_connector.Client, pair string, interval []Interval) []StrategyResult {
	var allResult []StrategyResult
	var wg sync.WaitGroup
	resultsChan := make(chan StrategyResult, 32*16) // Buffered channel

	// Step 1: Launch a collector goroutine
	go func() {
		for result := range resultsChan {
			allResult = append(allResult, result)
		}
	}()

	for i := 4; i < 15; i++ {
		for j := 10; j < 40; j += 2 {
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

	return allResult

}

// find average ratio

func FetchReports(client *binance_connector.Client) error {
	reports, err := OpenReports()
	if err != nil {
		return err
	}

	for _, r := range reports {
		result := ParralelFindBestMAParams(client, r, []Interval{m5, m15, m30, h1})
		err = AppendToReport(r, result)
		if err != nil {
			return err
		}
		time.Sleep(3 * time.Minute)
	}
	return nil
}

func FindAvgRatioPerParams() error {
	allReports, err := GetAllReports()
	if err != nil {
		return err
	}
	sumMap := make(map[string]float64)
	countMap := make(map[string]float64)
	avgMap := make(map[string]float64)
	for _, r := range allReports {
		str_param := fmt.Sprintf("%v_%v", r.Strategy.Main.Params[SMA_short], r.Strategy.Main.Params[SMA_long])
		countMap[str_param] += 1
		sumMap[str_param] += r.Ratio

	}
	for k, v := range sumMap {
		avgMap[k] = v / countMap[k]
	}
	type average struct {
		params string
		avg    float64
	}
	var avg []average
	for k, v := range avgMap {
		avg = append(avg, average{params: k, avg: v})
	}

	slices.SortFunc(avg, func(a, b average) int {
		if a.avg > b.avg {
			return 1
		}
		if b.avg > a.avg {
			return -1
		}
		return 0
	})
	fmt.Println(avg)
	return nil
}
