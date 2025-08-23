package main

import (
	"fmt"
	"log"
	"os"
	"path"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

func main() {
	// Load .env file
	curr, err := os.Getwd()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	err = godotenv.Load(path.Join(curr, "data", ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var apiKey = os.Getenv("API_KEY")
	var secretKey = os.Getenv("SECRET_KEY")

	client := binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")
	_ = client

	Run(client)

	// Get API credentials from environment variable

}

// creer des commandes pour backtester

func Test(client *binance_connector.Client) {

	var testStrat Wrapper = Wrapper{
		Asset:     "ETHUSDC",
		Amount:    0.002,
		Intervals: []Interval{m1, m5, m15, m30, h1},
		Main: Signal{
			Name:   "EMA",
			Type:   "Moving Average",
			Params: make(map[Indicator]int),
		},
	}

	testStrat.Main.Params[SMA_short] = 9
	testStrat.Main.Params[SMA_long] = 15
	testStrat.Main.Params[SMA_super_long] = 200
	r, err := testStrat.Run(client)
	if err != nil {
		fmt.Println(err)
	}
	r.SaveTradeResult(testStrat.Intervals[0])
}

func Run(client *binance_connector.Client) {
	var testStrat Wrapper = Wrapper{
		Asset:     "ETHUSDC",
		Amount:    0.002,
		Intervals: []Interval{m5, m15, m30, h1},
		Main: Signal{
			Name:   "EMA",
			Type:   "Moving Average",
			Params: make(map[Indicator]int),
		},
	}

	testStrat.Main.Params[SMA_short] = 9
	testStrat.Main.Params[SMA_long] = 15
	testStrat.Main.Params[SMA_super_long] = 200

	r, err := testStrat.Test(client)
	if err != nil {
		fmt.Println(err)
	}
	err = r.SaveTradeResult(testStrat.Intervals[0])
	if err != nil {
		fmt.Println(err)
	}
}
