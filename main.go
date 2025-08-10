package main

import (
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
	_ = client
	Processor(client)
}

func Processor(client *binance_connector.Client) {
	klines := IndicatorstoKlines(GetKlines(client, "ETHUSDC", "1h", 1000), 9, 25, 14)

	_ = klines
}

/*


	strats := []StrategyStat{
		{
			Asset:     "HBARUSDC",
			StratName: "RSI2575",
			Interval:  "1h",
			Ratio:     0,
		}, {
			Asset:     "ETHUSDC",
			StratName: "RSI2575",
			Interval:  "1h",
			Ratio:     0,
		}, {
			Asset:     "XRPUSDC",
			StratName: "RSI2575",
			Interval:  "5m",
			Ratio:     0,
		},
	}
	for _, strat := range strats {
		strat.SMATest(client)
		fmt.Println(strat.Ratio)
	}

*/
