package main

import (
	"fmt"
	"strconv"
)

func OverSuperLong(kline *Klines, i int) bool {
	f_close, err := strconv.ParseFloat(kline.Array[i].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f_close > kline.Indicators[SMA_super_long][i]
}
