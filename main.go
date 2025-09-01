package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	"github.com/Stupnikjs/binance_go/pkg/klines"
	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var apiKey = os.Getenv("API_KEY")
	var secretKey = os.Getenv("SECRET_KEY")

	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")
	_ = client

	if err != nil {
		fmt.Println(err)
	}

	// Get API credentials from environment variable
	k, err := klines.LoadKlinesFromFile(path.Join("data", string(klines.Interv[1]), strings.ToLower("BTCUSDC")))
	if err != nil {
		fmt.Println(err)
	}
	indic := []klines.Indicator{
		{Name: "RSI", Interval: klines.Interv[1], Type: "Price", Calculator: analysis.RSIcalc, Param: 14},
	}
	featured := klines.BuildFeaturedKlinesArray(k, indic)
	fmt.Println(klines.FeaturedKlinesToString(featured[0]))
}

func SaveLastKlines(client *binance_connector.Client, intervals []klines.Interval) error {
	for _, i := range PAIRS {
		err := klines.AppendNewData(client, i, intervals)
		if err != nil {
			return err
		}

	}
	return nil
}
