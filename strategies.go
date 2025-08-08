package main

import (
	"fmt"
	"log"

	binance_connector "github.com/binance/binance-connector-go"
)

type Bot struct {
	client     *binance_connector.Client
	strategies []*Strategy
	// ...
}

func (b *Bot) Run() {
	var savedTrade []Strategy

	for len(savedTrade) <= 1 {
		for _, s := range b.strategies {

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
	}

}
func RSIbuyCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, -25)
}
func RSIsellCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, 75)
}
