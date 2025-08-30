package klines

import (
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
	"reflect"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

// AppendToFile opens a file in append mode and encodes the new data to the end.
// This is more efficient than reading the entire file, appending, and then saving.
// check time continuity
func AppendToFile(data []*binance_connector.KlinesResponse, pair string, interval Interval) error {
	derefData := DeRefKlinesArray(data)
	path := path.Join("data", strings.ToLower(string(interval)), pair)

	// Load existing data from the file first.
	// This function handles opening, decoding, and closing the file.
	klines, err := LoadKlinesFromFile(path)
	if os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return fmt.Errorf("could not open file for writing: %w", err)
		}
		defer file.Close()

		// Create a single gob encoder for the entire process.
		encoder := gob.NewEncoder(file)

		// Encode the combined data in one go.
		if err := encoder.Encode(derefData); err != nil {
			return fmt.Errorf("could not encode data: %w", err)
		}

		return nil

	}

	if err != nil {
		return err
	}
	// Combine the existing data with the new data.
	// We're performing the data overlap and gap checks here on the combined data.
	if IsDataOverlap(klines, data) {
		data, err = SliceOverLaping(klines, data)
		if err != nil {
			return err
		}
	}

	if IsThereDataGap(klines, data) {
		return fmt.Errorf("there is a data gap")
	}

	combinedData := append(klines, data...)
	derefCombined := DeRefKlinesArray(combinedData)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("could not open file for writing: %w", err)
	}
	defer file.Close()

	// Create a single gob encoder for the entire process.
	encoder := gob.NewEncoder(file)

	// Encode the combined data in one go.
	if err := encoder.Encode(derefCombined); err != nil {
		return fmt.Errorf("could not encode data: %w", err)
	}

	return nil
}

// loadKlinesFromFile has been updated to read multiple gob-encoded objects
// from the file stream until it reaches the end of the file (io.EOF).
func LoadKlinesFromFile(filename string) ([]*binance_connector.KlinesResponse, error) {

	file, err := os.Open(filename)
	if os.IsNotExist(err) {
		return []*binance_connector.KlinesResponse{}, err
	}
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var allData []*binance_connector.KlinesResponse
	decoder := gob.NewDecoder(file)

	// Loop to decode all objects from the gob stream until an error or EOF is encountered.
	for {
		var data []*binance_connector.KlinesResponse
		if err := decoder.Decode(data); err != nil {
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
		return []*binance_connector.KlinesResponse{}, nil // All new data is already in old.
	}

	return new[index:], nil

}

func AppendNewData(client *binance_connector.Client, pair string, intervals []Interval) error {
	klines, err := FetchKlines(client, pair, intervals)
	if err != nil {
		return err
	}
	err = AppendToFile(klines, pair, intervals[0])
	if err != nil {
		return err
	}
	// find a better way to test
	filename := path.Join("data", string(intervals[0]), strings.ToLower(pair))
	newklines, err := LoadKlinesFromFile(filename)
	if err != nil {
		return err
	}
	if !reflect.DeepEqual(newklines[len(newklines)-1], klines[len(klines)-1]) {
		return fmt.Errorf("append to file not working ")
	}
	return nil
}

// to CSV
