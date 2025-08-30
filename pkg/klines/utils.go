package klines

import (
	"path"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

func DeRefKlinesArray(k []*binance_connector.KlinesResponse) []binance_connector.KlinesResponse {
	var arr []binance_connector.KlinesResponse
	for _, k := range k {
		arr = append(arr, *k)
	}
	return arr
}

func FileName(pair string, intervals []Interval) string {
	return path.Join("data", string(intervals[0]), strings.ToLower(pair))
}

func GetFileLen(pair string, intervals []Interval) int {
	klines, _ := LoadKlinesFromFile(FileName(pair, intervals))
	return len(klines)
}
