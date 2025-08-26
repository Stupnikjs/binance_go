package klines

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

// NOT WORKING
// append price to the file based on the interval and pair
func SaveKline(pair string, interval Interval) error {
	if !strings.HasSuffix(pair, "USDC") {
		return fmt.Errorf("give a pair with USDC")
	}
	var klineData []binance_connector.KlinesResponse
	_ = klineData
	return nil
}

func AppendKlineToFile(kline Klines, pair string, interval Interval) error {
	fileName := strings.ToLower(pair)
	var data []binance_connector.KlinesResponse
	for _, k := range kline.Array {
		data = append(data, *k)
	}
	return AppendToFile(data, fileName, interval)
}

// AppendToFile opens a file in append mode and encodes the new data to the end.
// This is more efficient than reading the entire file, appending, and then saving.
// check time continuity
func AppendToFile(data []binance_connector.KlinesResponse, filename string, interval Interval) error {
	// os.O_APPEND ensures we write to the end of the file.
	// os.O_CREATE creates the file if it doesn't exist, which is a good practice.
	// os.O_WRONLY is for write-only mode.
	path := path.Join("data", strings.ToLower(string(interval)), filename)
	file, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for appending: %w", err)
	}
	defer file.Close()

	klines, err := LoadKlinesFromFile(path, interval)
	if err != nil {
		return err
	}
	encoder := gob.NewEncoder(file)
	if IsDataOverlap(klines, data) {
		data, err := SliceOverLaping(klines, data)
		if err != nil {
			return err
		}
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("could not encode data: %w", err)

		}
		return nil
	} else {
		if err := encoder.Encode(data); err != nil {
			return fmt.Errorf("could not encode data: %w", err)

		}
	}
	// fix to only appen new data
	return fmt.Errorf("wait some time because data overlaps ")

}

// loadKlinesFromFile has been updated to read multiple gob-encoded objects
// from the file stream until it reaches the end of the file (io.EOF).
func LoadKlinesFromFile(filename string, interval Interval) ([]binance_connector.KlinesResponse, error) {

	path := path.Join("data", string(interval), filename)
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var allData []binance_connector.KlinesResponse
	decoder := gob.NewDecoder(file)

	// Loop to decode all objects from the gob stream until an error or EOF is encountered.
	for {
		var data []binance_connector.KlinesResponse
		if err := decoder.Decode(&data); err != nil {
			// If we reach the end of the file, break the loop.
			if err == io.EOF {
				break
			}
			// For any other error, return the error.
			return nil, fmt.Errorf("could not decode data: %w", err)
		}
		// Append the newly decoded data to the slice of all data.
		allData = append(allData, data...)
	}

	return allData, nil
}

func IsDataOverlap(old []binance_connector.KlinesResponse, new []binance_connector.KlinesResponse) bool {
	lastOld := old[len(old)-1]
	firstNew := new[0]

	if firstNew.CloseTime <= lastOld.CloseTime {
		return true
	} else {
		return false
	}
}

func SliceOverLaping(old []binance_connector.KlinesResponse, new []binance_connector.KlinesResponse) ([]binance_connector.KlinesResponse, error) {

	if !IsDataOverlap(old, new) {
		return nil, fmt.Errorf(" data isnt overlaping ")
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
		return []binance_connector.KlinesResponse{}, nil // All new data is already in old.
	}

	return new[index:], nil

}
