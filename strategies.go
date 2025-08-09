package main

import (
	binance_connector "github.com/binance/binance-connector-go"
)

func RSIbuyCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, -45)
}
func RSIsellCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, 65)
}

func EMAcrossOver(client *binance_connector.Client, pair string, EMAcoefShort int, EMAcoefLong int, over bool) bool {
	klines := GetKlines(client, pair, "1m", EMAcoefLong*2)
	closing := CloseFromKlines(klines)
	EMAshort := ExponetialMovingAverage(closing, EMAcoefShort)
	EMALong := ExponetialMovingAverage(closing, EMAcoefLong)

	bullishCross := EMAshort[len(EMAshort)-2] < EMALong[len(EMALong)-2] && EMAshort[len(EMAshort)-1] > EMALong[len(EMALong)-1]
	bearishCross := EMAshort[len(EMAshort)-2] > EMALong[len(EMALong)-2] && EMAshort[len(EMAshort)-1] < EMALong[len(EMALong)-1]
	if over {
		if bullishCross {
			return true
		}
	}
	if !over {
		if bearishCross {
			return true
		}
	}
	return false

}
