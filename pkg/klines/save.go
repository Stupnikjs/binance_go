package klines

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

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
	fileName := fmt.Sprintf("%s_%s", pair, string(interval))
	var data []binance_connector.KlinesResponse
	for _, k := range kline.Array {
		data = append(data, *k)
	}
	return AppendToFile(data, fileName)
}

// AppendToFile opens a file in append mode and encodes the new data to the end.
// This is more efficient than reading the entire file, appending, and then saving.
func AppendToFile(data []binance_connector.KlinesResponse, filename string) error {
	// os.O_APPEND ensures we write to the end of the file.
	// os.O_CREATE creates the file if it doesn't exist, which is a good practice.
	// os.O_WRONLY is for write-only mode.
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for appending: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("could not encode data: %w", err)
	}
	return nil
}

// loadKlinesFromFile has been updated to read multiple gob-encoded objects
// from the file stream until it reaches the end of the file (io.EOF).
func LoadKlinesFromFile(filename string) ([]binance_connector.KlinesResponse, error) {
	file, err := os.Open(filename)
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
