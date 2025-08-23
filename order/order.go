package order

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"time"

	binance_connector "github.com/binance/binance-connector-go"
)

type CreateOrderResponse struct {
	Symbol        string `json:"symbol"`
	OrderID       int64  `json:"orderId"`
	ClientOrderID string `json:"clientOrderId"`
	TransactTime  int64  `json:"transactTime"`
	Price         string `json:"price"`
	OrigQty       string `json:"origQty"`
	ExecutedQty   string `json:"executedQty"`
	Status        string `json:"status"`
	TimeInForce   string `json:"timeInForce"`
	Type          string `json:"type"`
	Side          string `json:"side"`
	Fills         []Fill `json:"fills"`
}

// Fill represents a single fill of the order.
type Fill struct {
	Price           string `json:"price"`
	Qty             string `json:"qty"`
	Commission      string `json:"commission"`
	CommissionAsset string `json:"commissionAsset"`
}

func ParseResponse(response interface{}) (*CreateOrderResponse, error) {
	var orderResponse CreateOrderResponse
	jsonBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(jsonBytes, &orderResponse)
	if err != nil {
		return nil, err
	}
	if len(orderResponse.Fills) == 0 {
		return nil, fmt.Errorf(" received empty response ")
	}
	return &orderResponse, nil
}

func BuildOrder(client *binance_connector.Client, orderType string, pair string, amount float64) (interface{}, error) {
	return client.NewCreateOrderService().
		Symbol(pair).
		Side(orderType).
		Type("MARKET").
		Quantity(amount).
		Do(context.Background())

}
func BuildStopLoss(client *binance_connector.Client, price float64, pair string, amount float64) error {
	order, err := BuildOrder(client, "STOPLOSS", pair, amount)
	if err != nil {
		return err
	}
	orderResp, err := ParseResponse(order)
	if err != nil {
		return err
	}

	fmt.Println(orderResp) // test
	return err
}

// store trade id from api
func TimeStampToDateString(stamp int) string {
	ts1 := int64(stamp)

	// Convert milliseconds to seconds and nanoseconds
	// The time.Unix() function takes seconds and nanoseconds
	seconds1 := ts1 / 1000
	nanoseconds1 := (ts1 % 1000) * 1000000

	// Create time.Time objects
	t1 := time.Unix(seconds1, nanoseconds1)
	return t1.Local().String()

}

func GetAssetBalance(client *binance_connector.Client, asset string) (float64, error) {

	account, err := client.NewGetAccountService().Do(context.Background())
	for i := range account.Balances {
		if asset == account.Balances[i].Asset {
			amount, err := strconv.ParseFloat(account.Balances[i].Free, 64)
			return amount, err
		}
	}
	return 0, err
}

func PrintUSDCBalance(client *binance_connector.Client) {
	usdc, err := GetAssetBalance(client, "USDC")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("USDC: %f \n", usdc)
}
