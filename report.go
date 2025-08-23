package main

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path"
	"strings"
)

func (r Result) SaveTradeResult(interval Interval) error {
	fileName := fmt.Sprintf("%s.json", strings.ToLower(r.Pair))
	path := path.Join("data", "trades", string(interval), fileName)
	file, err := os.Open(path)

	if err != nil {
		if err == os.ErrNotExist {
			file, err = os.Create(path)
			if err != nil {
				return err
			}

		}

	}
	defer file.Close()
	bytes, err := io.ReadAll(file)
	if err != nil {
		return err
	}

	var oldResult []Result
	err = json.Unmarshal(bytes, &oldResult)
	if err != nil {
		return err
	}

	if len(oldResult) <= 0 {
		oldResult = append(oldResult, r)
		bytes, err = json.Marshal(oldResult)
		if err != nil {
			return err
		}
		err = os.WriteFile(path, bytes, 0644)
		return err
	}
	if int64(oldResult[len(oldResult)-1].EndStamp) > int64(r.StartStamp) {
		return fmt.Errorf("result overlap wait for %d ", int64(oldResult[len(oldResult)-1].EndStamp)-int64(r.StartStamp))
	}

	oldResult = append(oldResult, r)
	file.Close()
	bytes, err = json.Marshal(oldResult)
	if err != nil {
		return err
	}
	err = os.WriteFile(path, bytes, 0644)
	fmt.Println(r)
	return err

}
