package main

// save all trades []Trade => .txt .xls

import (
	"encoding/json"
	"os"
)

func CreateHistoryFile(fileName string) (*os.File, error) {
	// os.O_CREATE: Create the file if it doesn't exist.
	// os.O_RDWR: Open the file for reading and writing.
	// os.O_APPEND: Append to the end of the file when writing.
	return os.OpenFile(fileName, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0666)
}

func SaveTrade(s Trader) error {
	file, err := CreateHistoryFile("history.txt")

	if err != nil {
		return err
	}
	jsonBytes, err := json.Marshal(s)
	file.Write(jsonBytes)
	defer file.Close()
	return err
}
