package main

import (
	"context"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

func GetAssetBalance(client *binance_connector.Client, asset string) (error, float64) {

	account, err := client.NewGetAccountService().Do(context.Background())
	for i := range account.Balances {
		if asset == account.Balances[i].Asset {
			amount, err := strconv.ParseFloat(account.Balances[i].Free, 64)
			return err, amount
		}
	}
	return err, 0
}
