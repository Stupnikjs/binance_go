package klines

import (
	"encoding/csv"
	"encoding/gob"
	"fmt"
	"os"
	"path"
	"strings"

	binance_connector "github.com/binance/binance-connector-go"
)

func AppendToFile(data []*binance_connector.KlinesResponse, pair string, interval Interval) (int, error) {
	derefData := DeRefKlinesArray(data)
	path := path.Join("data", strings.ToLower(string(interval)), pair)

	klines, err := LoadKlinesFromFile(path)

	// case where file is created yet
	if os.IsNotExist(err) {
		file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			return 0, fmt.Errorf("could not open file for writing: %w", err)
		}
		defer file.Close()

		// Create a single gob encoder for the entire process.
		encoder := gob.NewEncoder(file)

		// Encode the combined data in one go.
		if err := encoder.Encode(derefData); err != nil {
			return 0, fmt.Errorf("could not encode data: %w", err)
		}

		return len(derefData), nil

	}

	if err != nil {
		return 0, err
	}

	if IsDataOverlap(klines, data) {
		data, err = SliceOverLaping(klines, data)
		if err != nil {
			return 0, err
		}
	}

	if IsThereDataGap(klines, data) {
		return 0, fmt.Errorf("there is a data gap")
	}

	combinedData := append(klines, data...)
	derefCombined := DeRefKlinesArray(combinedData)

	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return 0, fmt.Errorf("could not open file for writing: %w", err)
	}
	defer file.Close()

	// Create a single gob encoder for the entire process.
	encoder := gob.NewEncoder(file)

	// Encode the combined data in one go.
	if err := encoder.Encode(derefCombined); err != nil {
		return 0, fmt.Errorf("could not encode data: %w", err)
	}

	return len(derefCombined), nil
}

func AppendNewData(client *binance_connector.Client, pair string, intervals []Interval) error {
	klines, err := FetchKlines(client, pair, intervals)
	if err != nil {
		return err
	}
	length, err := AppendToFile(klines, pair, intervals[0])
	if err != nil {
		return err
	}
	fmt.Printf("%s file has %d lines \n", pair, length)
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
	file, err := os.Create(filename + ".csv")
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
