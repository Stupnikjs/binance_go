package main

import binance_connector "github.com/binance/binance-connector-go"

func RSIbuyCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, -25)
}
func RSIsellCondition14(client *binance_connector.Client, pair string, interval string) bool {
	return RSIstrat(client, pair, interval, 14, 75)
}
