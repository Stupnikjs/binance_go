package klines

import (
	"encoding/gob"
	"fmt"
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

func SaveKlineToFile(kline Klines, pair string, interval Interval) error {
	fileName := fmt.Sprintf("%s_%s", pair, string(interval))
	var data []binance_connector.KlinesResponse
	for _, k := range kline.Array {
		data = append(data, *k)
	}
	return SaveToFile(data, fileName)
}

func SaveToFile(data []binance_connector.KlinesResponse, filename string) error {
	file, err := os.Create(filename)
	if err != nil {
		if err != os.ErrExist {
			// call the append method
			return nil
		}
		return fmt.Errorf("could not create file: %w", err)
	}
	defer file.Close()

	encoder := gob.NewEncoder(file)
	if err := encoder.Encode(data); err != nil {
		return fmt.Errorf("could not encode data: %w", err)
	}
	return nil
}

func loadKlinesFromFile(filename string) ([]binance_connector.KlinesResponse, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("could not open file: %w", err)
	}
	defer file.Close()

	var data []binance_connector.KlinesResponse
	decoder := gob.NewDecoder(file)
	if err := decoder.Decode(&data); err != nil {
		return nil, fmt.Errorf("could not decode data: %w", err)
	}

	return data, nil
}
