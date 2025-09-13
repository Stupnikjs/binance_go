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
	var refData []*binance_connector.KlinesResponse
	for _, k := range allData {
		refData = append(refData, &k)
	}
	fmt.Printf("%d klines loaded from %s \n ", len(refData), filename)
	return refData, nil
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
