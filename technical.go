package main

import (
	"context"
	"fmt"
	"math"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

const (
	RSI     Indicator = "RSI"
	VOL     Indicator = "VOL"
	RSI_15m Indicator = "RSI_15m"
	RSI_30m Indicator = "RSI_30m"
	RSI_1h  Indicator = "RSI_1h"
	RSI_2h  Indicator = "RSI_2h"
	RSI_4h  Indicator = "RSI_4h"

	SMA_short Indicator = "SMA_short"
	EMA_short Indicator = "EMA_short"
	SMA_long  Indicator = "SMA_long"
	EMA_long  Indicator = "EMA_long"
)

type Interval string

const (
	m1  Interval = "1m"
	m5  Interval = "5m"
	m15 Interval = "15m"
	m30 Interval = "30m"
	h1  Interval = "1h"
	h2  Interval = "2h"
	h4  Interval = "4h"
)

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type Indicators map[Indicator][]float64
type Klines struct {
	Interval   Interval
	Array      []*binance_connector.KlinesResponse
	Indicators Indicators
}

type IndicatorsParams struct {
	short_period_MA int
	long_period_MA  int
	RSI_coef        int
}

// Kline get upper interval
// query upper Intervals Coefs ex RSI_1h ..

// volume weighted Average Price

func BuildKlinesArr(client *binance_connector.Client, pair string, Interval []Interval) []*Klines {
	klinesArr := []*Klines{}
	for _, i := range Interval {
		klines, err := client.NewKlinesService().
			Symbol(pair).
			Interval(string(i)).
			Limit(1000).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		kl := Klines{
			Array:      klines,
			Interval:   i,
			Indicators: make(Indicators),
		}
		klinesArr = append(klinesArr, &kl)
	}
	return klinesArr
}

// refactor
func CloseFromKlines(klines []*binance_connector.KlinesResponse) []float64 {
	closingPrices := make([]float64, len(klines))
	for i, kline := range klines {
		f_close, err := strconv.ParseFloat(kline.Close, 64)
		if err != nil {
			fmt.Println(err)
		}
		closingPrices[i] = f_close

	}
	return closingPrices
}

// error checking
func IndicatorstoKlines(client *binance_connector.Client, pair string, intervals []Interval, params IndicatorsParams) []*Klines {
	klinesArr := BuildKlinesArr(client, pair, intervals)
	ProcessKlines(klinesArr, params)
	return klinesArr
}

func ProcessKlines(klines []*Klines, params IndicatorsParams) []*Klines {
	for _, k := range klines {
		// caclulate RSI EMA SMA
		close := CloseFromKlines(k.Array)
		RSI_arr := RSIcalc(close, params.RSI_coef)
		SMA_short_arr := SMAcalc(close, params.short_period_MA)
		SMA_long_arr := SMAcalc(close, params.long_period_MA)
		EMA_short_arr := EMAcalc(close, params.short_period_MA)
		EMA_long_arr := EMAcalc(close, params.long_period_MA)
		k.Indicators[RSI] = RSI_arr
		k.Indicators[SMA_short] = SMA_short_arr
		k.Indicators[SMA_long] = SMA_long_arr
		k.Indicators[EMA_short] = EMA_short_arr
		k.Indicators[EMA_long] = EMA_long_arr
	}
	return klines
}

func SMAcalc(closingPrices []float64, period int) []float64 {
	var SMA []float64
	closingPriceSlice := closingPrices
	if len(closingPrices) < period {
		return SMA
	}
	for i := period - 1; i < len(closingPriceSlice); i++ {
		var sma float64
		slice := closingPriceSlice[i-period+1 : i+1]
		for _, n := range slice {
			sma += n

		}
		SMA = append(SMA, sma/float64(period))

	}
	return SMA
}

func RSIcalc(prices []float64, period int) []float64 {
	// Le RSI ne peut pas être calculé si le nombre de prix est inférieur à la période.
	if len(prices) <= period {
		return nil
	}

	// Initialiser les slices pour les gains, les pertes et le RSI
	gains := make([]float64, len(prices))
	losses := make([]float64, len(prices))
	rsi := make([]float64, len(prices))

	// Étape 1 & 2 : Calculer les changements de prix, les gains et les pertes
	for i := 1; i < len(prices); i++ {
		change := prices[i] - prices[i-1]
		if change > 0 {
			gains[i] = change
			losses[i] = 0
		} else {
			gains[i] = 0
			losses[i] = math.Abs(change)
		}
	}

	// Étape 3 : Calculer la première moyenne de gain et de perte (moyenne simple)
	var avgGain float64
	var avgLoss float64
	for i := 1; i <= period; i++ {
		avgGain += gains[i]
		avgLoss += losses[i]
	}
	avgGain /= float64(period)
	avgLoss /= float64(period)

	// Étape 4 : Calculer le premier RS et le premier RSIcalc
	// Le premier RSIcalc est stocké à l'index `period`
	if avgLoss == 0 {
		rsi[period] = 100 // Pour éviter la division par zéro
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// Étape 5 : Calculer les RSIcalc suivants avec la méthode de lissage
	for i := period + 1; i < len(prices); i++ {
		avgGain = ((avgGain * float64(period-1)) + gains[i]) / float64(period)
		avgLoss = ((avgLoss * float64(period-1)) + losses[i]) / float64(period)

		if avgLoss == 0 {
			rsi[i] = 100
		} else {
			rs := avgGain / avgLoss
			rsi[i] = 100 - (100 / (1 + rs))
		}
	}

	// Les premières (période) valeurs sont 0, on pourrait retourner une slice plus courte si désiré
	return rsi
}

func EMAcalc(closingPrices []float64, period int) []float64 {

	if len(closingPrices) <= period {
		return SMAcalc(closingPrices, period)
	}

	firstSMA := SMAcalc(closingPrices, period)[0]
	EMA := []float64{}

	EMA = append(EMA, firstSMA)
	EMAcoef := 2.0 / float64(period+1)
	prevEMA := firstSMA
	for i := period; i < len(closingPrices); i++ {
		nextEMA := closingPrices[i]*float64(EMAcoef) + prevEMA*(1-EMAcoef)
		prevEMA = nextEMA
		EMA = append(EMA, nextEMA)

	}
	return EMA
}
