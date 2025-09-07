package klines

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"io"
	"os"
	"path"
	"strconv"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

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

	fmt.Printf("%d klines loaded \n", len(refData))
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

func AppendNewData(client *binance_connector.Client, pair string, intervals []Interval) error {
	klines, err := FetchKlines(client, pair, intervals)
	if err != nil {
		return err
	}
	err = AppendToFile(klines, pair, intervals[0])
	if err != nil {
		return err
	}

	return nil
}

func CheckWholeHasNoTimeGap(pair string, interval Interval) error {
	klines, err := LoadKlinesFromFile(path.Join("data", string(interval), strings.ToLower(pair)))
	if err != nil {
		return err
	}
	prevGap := klines[1].CloseTime - klines[0].CloseTime
	for i := 1; i < len(klines)-1; i++ {

		gap := klines[i+1].CloseTime - klines[i].CloseTime
		if gap != prevGap {
			fmt.Println(klines[i+1].CloseTime, klines[i].CloseTime)
			return fmt.Errorf("gap supposed to be constant")
		}
		prevGap = gap

	}
	return nil
}

// to CSV

func FeaturedKlinesToCSV(filename string, data []FeaturedKlines) error {
	// Créer le fichier.
	file, err := os.Create(filename)
	if err != nil {
		return fmt.Errorf("impossible de créer le fichier %s : %w", filename, err)
	}
	defer file.Close()

	// Initialiser un nouveau writer CSV.
	writer := csv.NewWriter(file)
	defer writer.Flush()
	var headers []string
	// Écrire les en-têtes de colonnes.

	headers = append(headers, "price")
	headers = append(headers, "volume")
	headers = append(headers, "closing_time")

	for k := range data[0].FeaturesMap {
		headers = append(headers, k)
	}

	if err := writer.Write(headers); err != nil {
		return fmt.Errorf("impossible d'écrire les en-têtes : %w", err)
	}

	for _, k := range data {
		stringed := FeaturedKlinesToString(k)
		writer.Write(stringed)
	}

	return nil
}

func FeaturedKlinesToString(f FeaturedKlines) []string {
	var arr []string
	arr = append(arr, f.Close)
	arr = append(arr, f.Volume)
	strTime := strconv.Itoa(int(f.CloseTime))
	arr = append(arr, strTime)
	for _, v := range f.FeaturesMap {
		str := strconv.FormatFloat(v, 'f', 3, 64)
		arr = append(arr, str)
	}
	return arr
}
