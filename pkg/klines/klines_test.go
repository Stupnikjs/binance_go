package klines

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
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

/*
func TestCloseFromKlines(t *testing.T) {

	klines, err := LoadKlinesFromFile(GetFilePathName("BTCUSDC", Interv[1]))
	if err != nil {
		t.Error("err should be nil")
	}
	close := CloseFromKlines(klines)
	if len(close) != len(klines) {
		t.Error("close should be equal to klines Array")
	}

	f_close, err := strconv.ParseFloat(klines[0].Close, 64)

	if err != nil {
		t.Errorf("float conversion gave error %v", err)
	}
	if close[0] != f_close {
		t.Error("first value of array and close array are diffrent")
	}
}

*/

func TestLoadFromKlines(t *testing.T) {
	curr, _ := os.Getwd()
	path := filepath.Join(curr, "../../data/5m/BTCUSDC")
	k, err := LoadKlinesFromFile(path)
	// MARCHE PAS
	if err != nil {
		t.Errorf("err loading klines %v", err)
	}
	fmt.Println(k)

}
