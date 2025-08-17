package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"slices"
	"strings"
)

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
	fmt.Printf("appening %s to report .. \n", pair)
	filename := fmt.Sprintf("%s_report.json", strings.ToLower(pair))
	path := path.Join("reports_5m", filename)
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

func PrintTop5report(pairname string) error {
	results, err := ReadReport(pairname)
	if err != nil {
		return err

	}
	slices.SortFunc(results, func(a, b StrategyResult) int {
		if a.Ratio > b.Ratio {
			return -1
		}
		if a.Ratio < b.Ratio {
			return 1
		}
		return 0
	})

	for _, r := range results[:5] {
		fmt.Printf("Top 5 %s with ratio:%v params %v", pairname, r.Ratio, r.Strategy.Main)
	}
	return nil
}

func GetAllReports() ([]StrategyResult, error) {
	names, err := OpenReports()
	allResult := []StrategyResult{}
	if err != nil {
		return nil, err
	}
	for _, n := range names {
		result, err := ReadReport(n)
		if err != nil {
			return nil, err
		}
		allResult = append(allResult, result...)
	}
	return allResult, nil
}

func OpenReports() ([]string, error) {
	pairs := []string{}
	entries, err := os.ReadDir("reports_5m")
	if err != nil {
		return nil, err
	}
	for _, e := range entries {
		pair := strings.Split(e.Name(), "_")

		if len(pair) > 1 {
			pairs = append(pairs, strings.ToUpper(pair[0]))
		} else {
			return nil, fmt.Errorf("report name doesnt matchs expetation")
		}

	}
	return pairs, nil
}

func ReadReport(pairname string) ([]StrategyResult, error) {
	fileName := fmt.Sprintf("%s_report.json", strings.ToLower(pairname))
	path := path.Join("reports_5m", fileName)
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
