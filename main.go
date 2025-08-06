package main

import (
	"fmt"
	"log"
	"os"
	"time"

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

	Processor(client)

}

func Processor(client *binance_connector.Client) {

	activeTrade := []Trade{}

	for {
		print(RSIunder25(client, "ETHUSDC", "1m", 40, 14))
		if true {
			trade := Trade{
				Asset:  "ETHUSDC",
				Amount: 0.002,
			}
			err := trade.Buy(client)
			if err != nil {
				fmt.Println(err)
			}
			log.Println(trade)
			activeTrade = append(activeTrade, trade)
		}

		time.Sleep(3 * time.Minute)
		// check for active Trades

		if len(activeTrade) > 0 {

			// check for sell condtions
			err := activeTrade[0].Sell(client)
			log.Println(activeTrade[0])
			if err != nil {
				fmt.Println(err)
			}
		}

	}

}

func TradingScript(client *binance_connector.Client) {

	// fonction qui va checker les indicateurs et conditions d'achats

	// qui retourne un []String avec les paires a acheter

	// loop sur []String pour creer les trades et executer les ordres

	// check sur les indicateurs et conditions de vente sur les trade en cours
	// cloture des trades en cours et vente si indicateurs valident

	// incrementation d'un []Trade de trade finis

	// sortie de la boucle ==> edition d'un fichier de sythese des trades

}
