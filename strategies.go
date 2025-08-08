package main

import (
	"fmt"
	"log"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

type Bot struct {
	client     *binance_connector.Client
	strategies []*Strategy
	// ...
}

func (b *Bot) Run() {
	var savedTrade []Strategy

	for {
		for _, s := range b.strategies {
			// bloquer l'achat sur une strategy
			if s.BuyCondition(b.client, s.Asset, "1m") {
				err := s.Buy(b.client)
				usdc_balance, _ := GetAssetBalance(b.client, "USDC")
				log.Println(usdc_balance)
				log.Println(s.Trade)
				if err != nil {
					fmt.Println(err)
				}

			}
			if s.Trade != nil {
				if s.SellCondition(b.client, s.Asset, "1m") {
					err := s.Sell(b.client)
					if err != nil {
						fmt.Println(err)
					}
					SaveTrade(*s)
					usdc_balance, _ := GetAssetBalance(b.client, "USDC")
					log.Println(usdc_balance)
					savedTrade = append(savedTrade, *s)

				}

			}

		}
		time.Sleep(40 * time.Second)
	}

}
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
