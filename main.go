package main

import (
	"fmt"
	"log"
	"os"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get API credentials from environment variables
	apiKey := os.Getenv("API_KEY")
	secretKey := os.Getenv("SECRET_KEY")

	// Initialize Binance client
	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")

	// Example: Get account information
	_ = client

	//Processor(client)
	for {
		fmt.Print(">:")
		cmd := Prompt()
		fmt.Println(cmd)
	}

}

func Processor(client *binance_connector.Client) {
	stratQueue := []Strategy{}
	savedTrade := []Strategy{}
	newStrat := Strategy{
		Asset:         "BTCUSDC",
		Amount:        0.001,
		BuyCondition:  RSIbuyCondition14,
		SellCondition: RSIsellCondition14,
	}
	stratQueue = append(stratQueue, newStrat)
	for len(savedTrade) <= 1 {
		if newStrat.BuyCondition(client, newStrat.Asset, "1m") {
			err := newStrat.Buy(client)
			usdc_balance, _ := GetAssetBalance(client, "USDC")
			log.Println(usdc_balance)
			log.Println(newStrat.Trade)
			if err != nil {
				fmt.Println(err)
			}

		}
		if newStrat.SellCondition(client, newStrat.Asset, "1m") {
			err := newStrat.Sell(client)
			if err != nil {
				fmt.Println(err)
			}
			SaveTrade(newStrat)
			usdc_balance, _ := GetAssetBalance(client, "USDC")
			log.Println(usdc_balance)
			savedTrade = append(savedTrade, newStrat)

		}

	}

}
