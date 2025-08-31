package main

import (
	"fmt"
	"log"
	"os"
	"path"
	"strings"

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

	for _, i := range PAIRS {
		k, _ := klines.LoadKlinesFromFile(path.Join("data", string(klines.Interv[1]), strings.ToLower(i)))
		Featured := klines.BuildFeaturedKlinesArray(k, klines.Indicators)
		fmt.Println(Featured[len(Featured)-3:])
		if err != nil {
			fmt.Println(err)

		}
	}

}
