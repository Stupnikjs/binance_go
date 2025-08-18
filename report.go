package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

func ReadReport(interval Interval) ([]StrategyResult, error) {

	//   read the report based on interval

	fileName := fmt.Sprintf("report_%s.json", string(interval))

	file, err := os.Open(fileName)
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

func AppendToReport(pair string, results []StrategyResult, interval Interval) error {
	path := fmt.Sprintf("report_%s.json", interval)
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

	if len(oldResult) <= 0 {
		bytes, err = json.Marshal(results)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, bytes, 0644)
		return err
	}
	fmt.Println(oldResult)
	diff := oldResult[len(oldResult)-1].EndStamp - results[0].StartStamp
	if diff > 0 {
		return fmt.Errorf("wait %v second for old period to be over", diff)
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
