package klines

import (
	"fmt"
	"path"
	"strings"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

func IntervalToTime(interval Interval) (time.Duration, error) {
	// The time.ParseDuration function can handle strings like "10m", "1h", "1m30s".
	// It returns a duration and an error.
	return time.ParseDuration(string(interval))
}

func FileName(pair string, intervals []Interval) string {
	return path.Join("data", string(intervals[0]), strings.ToLower(pair))
}

func DeRefKlinesArray(k []*binance_connector.KlinesResponse) []binance_connector.KlinesResponse {
	var arr []binance_connector.KlinesResponse
	for _, k := range k {
		arr = append(arr, *k)
	}
	return arr
}

func GetFileLen(pair string, intervals []Interval) int {
	klines, _ := LoadKlinesFromFile(FileName(pair, intervals))
	return len(klines)
}

func IsThereDataGap(old []*binance_connector.KlinesResponse, new []*binance_connector.KlinesResponse) bool {
	if len(old) < 1 || len(new) < 1 {
		return false
	}
	lastOld := old[len(old)-1]
	firstNew := new[0]
	regularGap, err := GetTimeGap(old)
	if err != nil {
		fmt.Println(err)
	}
	if firstNew.CloseTime-lastOld.CloseTime > regularGap {
		return true
	}
	return false
}

func GetTimeGap(kline []*binance_connector.KlinesResponse) (uint64, error) {
	if len(kline) >= 0 {
		return kline[1].CloseTime - kline[0].CloseTime, nil
	}
	return 0, fmt.Errorf("kline must be at least of len 2")
}

func SliceOverLaping(old []*binance_connector.KlinesResponse, new []*binance_connector.KlinesResponse) ([]*binance_connector.KlinesResponse, error) {
	if !IsDataOverlap(old, new) {
		return nil, fmt.Errorf("data isnt overlaping ")
	}
	lastOld := old[len(old)-1]
	var index int
	for i, n := range new {
		if n.CloseTime > lastOld.CloseTime {
			index = i
			break
		}
	}

	if index == 0 && new[0].CloseTime <= lastOld.CloseTime {
		return []*binance_connector.KlinesResponse{}, nil // All new data is already in old.
	}

	return new[index:], nil

}

func GetFilePathName(pair string, interval Interval) string {
	return path.Join("data", string(interval), pair)
}
