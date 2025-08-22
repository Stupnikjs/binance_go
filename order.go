package main

import (
	"context"
	"encoding/json"
	"fmt"
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

func (t *LiveTrader) ParseResponse(response interface{}) (*CreateOrderResponse, error) {
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

func (t *LiveTrader) BuildOrder(client *binance_connector.Client, orderType string) (interface{}, error) {
	return client.NewCreateOrderService().
		Symbol(t.Asset).
		Side(orderType).
		Type("MARKET").
		Quantity(t.Amount).
		Do(context.Background())

}

func (t *LiveTrader) BuildStopLoss(client *binance_connector.Client, price float64) error {
	order, err := t.BuildOrder(client, "STOPLOSS")
	if err != nil {
		return err
	}
	orderResp, err := t.ParseResponse(order)
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
