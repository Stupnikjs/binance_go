package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"slices"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

func FindBestMAParams(client *binance_connector.Client, pair string, interval []Interval) {
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
	for i := 4; i < 30; i++ {
		for j := 8; j < 60; j++ {
			if i >= j {
				continue
			} // i is < to j now
			s.Main.Params[SMA_short] = i
			s.Main.Params[SMA_long] = j
			m.Main.Params[SMA_short] = i
			m.Main.Params[SMA_long] = j
			r := s.StrategyTester(client)
			r.AppendToHistory()
			l := m.StrategyTester(client)
			l.AppendToHistory()
		}
	}

	// AnalyseReport(pair)
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
