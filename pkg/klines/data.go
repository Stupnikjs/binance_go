package klines

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"

	binance_connector "github.com/binance/binance-connector-go"
)

func LoadKlinesFromFile(filename string) ([]*binance_connector.KlinesResponse, error) {
	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return []*binance_connector.KlinesResponse{}, err
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allData []binance_connector.KlinesResponse
	decoder := gob.NewDecoder(file)

	// Loop to decode all objects from the gob stream until an error or EOF is encountered.
	for {
		var data []binance_connector.KlinesResponse
		if err := decoder.Decode(&data); err != nil {
			if err == io.EOF {
				break
			}
			return nil, fmt.Errorf("could not decode data: %w", err)
		}
		allData = append(allData, data...)
	}
	refData := ReRefKlines(allData)
	fmt.Printf("%d klines loaded from %s \n ", len(refData), filename)
	return refData, nil
}

func ReRefKlines(deRefed []binance_connector.KlinesResponse) []*binance_connector.KlinesResponse {
	var refData []*binance_connector.KlinesResponse
	for _, k := range deRefed {
		refData = append(refData, &k)
	}

	return refData
}

func IsDataOverlap(old []*binance_connector.KlinesResponse, new []*binance_connector.KlinesResponse) bool {
	if len(old) < 1 || len(new) < 1 {
		return false
	}
	lastOld := old[len(old)-1]
	firstNew := new[0]

	if firstNew.CloseTime <= lastOld.CloseTime {
		return true
	} else {
		return false
	}

}
