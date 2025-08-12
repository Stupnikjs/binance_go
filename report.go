package main

import (
	"encoding/json"
	"fmt"
	"os"
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

	// Overwrite the file with the new JSON data.
	// os.WriteFile handles opening, truncating, writing, and closing.
	err = os.WriteFile("report.json", finalBytes, 0644)
	if err != nil {
		return fmt.Errorf("error writing to file: %w", err)
	}

	return nil
}

func DataFromReport() {
	// trouver les meilleurs coefs par interval
}
