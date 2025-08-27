package klines

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"testing"

	binance_connector "github.com/binance/binance-connector-go"
	"github.com/joho/godotenv"
)

// TestMain is the entry point for the test suite.

var client *binance_connector.Client

func TestMain(m *testing.M) {
	// Call the setup function.
	err := godotenv.Load("../../.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	var apiKey = os.Getenv("API_KEY")
	var secretKey = os.Getenv("SECRET_KEY")

	client = binance_connector.NewClient(apiKey, secretKey, "https://testnet.binance.vision")

	if err != nil {
		fmt.Println(err)
	}
	// Run all tests.
	code := m.Run()

	// Call the teardown function.

	// Exit with the appropriate code.
	os.Exit(code)
}

func TestCloseFromKlines(t *testing.T) {

	klines := BuildKlineArrData(client, "BTCUSDC", Interv)
	close := CloseFromKlines(klines[0].Array)
	if len(close) != len(klines[0].Array) {
		t.Error("close should be equal to klines Array")
	}

	f_close, err := strconv.ParseFloat(klines[0].Array[0].Close, 64)

	if err != nil {
		t.Errorf("float conversion gave error %v", err)
	}
	if close[0] != f_close {
		t.Error("first value of array and close array are diffrent")
	}
}

func TestBuildKlinesArr(t *testing.T) {
	klinesArr := BuildKlinesArr(client, "BTCUSD", Interv)

	if len(klinesArr) != len(Interv) {
		t.Error("klines arr should be same len as Interval array ")
	}
	// check if data is continuous

}

func TestIndicatorToKlines(t *testing.T) {
	pair := "BTCUSDC"
	params := IndicatorsParams{
		Short_period_MA: 3,
		Long_period_MA:  10,
		Super_long_MA:   200,
		RSI_coef:        14,
	}
	_ = BuildKlineArrData(pair, Interv)
	_ = IndicatorstoKlines(client, pair, Interv, params)

}
