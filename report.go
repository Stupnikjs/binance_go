package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func ReadReport(pairname string, interval Interval) ([]StrategyResult, error) {

	//   read the report based on interval

	fileName := fmt.Sprintf("%s_report.json", strings.ToLower(pairname))
	var interval_string string
	switch interval {
	case m5:
		interval_string = "5m"
	case m15:
		interval_string = "15m"
	default:
		return nil, fmt.Errorf("interval must be provided")
	}
	path := path.Join("reports_"+interval_string, fileName)
	file, err := os.Open(path)
	if err != nil {
		return nil, err

	}
	bytes, err := io.ReadAll(file)
	if err != nil {
		return nil, err
	}
	var results []StrategyResult
	err = json.Unmarshal(bytes, &results)
	if err != nil {
		return nil, err

	}
	return results, nil
}

func (r *StrategyResult) AppendToHistory() error {
	// Read the existing file content.
	// os.ReadFile is cleaner than Open/ReadAll/Close for this task.
	bytes, err := os.ReadFile("report.json")
	if err != nil && !os.IsNotExist(err) {
		// Handle all errors except for the file not existing.
		return fmt.Errorf("error reading file: %w", err)
	}

	var storedResults []StrategyResult
	// If the file is not empty, unmarshal its content.
	if len(bytes) > 0 {
		// Pass the address of the slice to json.Unmarshal.
		// This is crucial for the function to be able to modify the slice.
		err = json.Unmarshal(bytes, &storedResults)
		if err != nil {
			return fmt.Errorf("error unmarshaling JSON: %w", err)
		}
	}

	// Append the new result to the slice.
	storedResults = append(storedResults, *r)

	// Marshal the updated slice back into JSON.
	finalBytes, err := json.Marshal(storedResults)
	if err != nil {
		return fmt.Errorf("error marshaling JSON: %w", err)
	}

	filename := fmt.Sprintf("%s_report.json", strings.ToLower(r.Strategy.Asset))

	// Overwrite the file with the new JSON data.
	// os.WriteFile handles opening, truncating, writing, and closing.
	err = os.WriteFile(filename, finalBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func SaveJsonResult(filename string, results []StrategyResult) {
	finalBytes, err := json.Marshal(results)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(filename, finalBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

}
func SaveJsonTrader(filename string, Traders []Trader) {
	finalBytes, err := json.Marshal(Traders)
	if err != nil {
		fmt.Println(err)
	}
	err = os.WriteFile(filename, finalBytes, 0644)
	if err != nil {
		fmt.Println(err)
	}

}

func AppendToReport(pair string, results []StrategyResult) error {
	path := "reports.json"
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}
	var oldResult []StrategyResult
	err = json.Unmarshal(bytes, &oldResult)
	if err != nil {
		return err
	}
	oldResult = append(oldResult, results...)
	file.Close()
	bytes, err = json.Marshal(oldResult)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	return err

}
