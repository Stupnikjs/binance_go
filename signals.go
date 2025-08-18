package main

import (
	"fmt"
	"strconv"
)

func OverSuperLong(kline *Klines, offset int, i int) bool {
	f_close, err := strconv.ParseFloat(kline.Array[offset].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return f_close > kline.Indicators[SMA_super_long][i]
}

func (t *Trader) CrossMA(klines *Klines, i int, index_super_long int, bigOverSmallPrev *bool, closed *[]Trader, test bool) bool {
	offset := i + index_super_long - 1
	closeOverMAsuperLong := OverSuperLong(klines, offset, i)
	bigOverSmall := klines.Indicators[SMA_short][i] < klines.Indicators[SMA_long][i]

	if !bigOverSmall && *bigOverSmallPrev && closeOverMAsuperLong {
		if test {
			t.BuyTest(klines, offset)
		} else {
			t.Buy(t.Client)
		}

	}
	if bigOverSmall && !*bigOverSmallPrev && t.Buy_time != 0 {
		if test {
			t.SellTest(klines, offset, closed)
		} else {
			t.Sell(t.Client)
		}
	}
	return bigOverSmall

}
