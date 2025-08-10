package main

import (
	"context"
	"fmt"
	"math"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

type Kline struct {
	Kline_binance *binance_connector.KlinesResponse
	Indicators    map[string]float64
}

func GetKlines(client *binance_connector.Client, pair string, interval string, limit int) []*binance_connector.KlinesResponse {
	klines, err := client.NewKlinesService().
		Symbol(pair).
		Interval(interval).
		Limit(limit).
		Do(context.Background())
	if err != nil {
		fmt.Println(err)
	}

	return klines
}

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

func IndicatorstoKlines(klines []*binance_connector.KlinesResponse, smallPeriod int, bigPeriod int, RSIcoef int) []Kline {
	klineArr := []Kline{}
	close := CloseFromKlines(klines)
	small_SMA := SMA(close, smallPeriod)
	big_SMA := SMA(close, bigPeriod)
	RSI := RSI(close, RSIcoef)
	fmt.Println(len(RSI), len(close))

	small_sma_field := fmt.Sprintf("sma_%d", smallPeriod)
	big_sma_field := fmt.Sprintf("sma_%d", bigPeriod)
	small_period_index := -smallPeriod
	big_period_index := -bigPeriod
	for _, k := range klines {
		small_period_index += 1
		big_period_index += 1
		kl := Kline{
			Kline_binance: k,
			Indicators:    make(map[string]float64),
		}

		if small_period_index >= 0 {
			kl.Indicators[small_sma_field] = small_SMA[small_period_index]
			if big_period_index >= 0 {
				kl.Indicators[big_sma_field] = big_SMA[big_period_index]
			}

		}

		klineArr = append(klineArr, kl)

	}
	return klineArr
}

func SMA(closingPrices []float64, period int) []float64 {
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

func IsMACrossOver(closingPrices []float64, sma []float64) bool {
	if len(closingPrices) < 2 || len(sma) < 2 {
		return false
	}

	if closingPrices[len(closingPrices)-2] > sma[len(sma)-2] {
		// price already above sma
		return false
	}
	if closingPrices[len(closingPrices)-2] < sma[len(sma)-2] &&
		closingPrices[len(closingPrices)-1] >= sma[len(sma)-1] {
		return true
	}
	return false
}

func PairCrossOverSignal(client *binance_connector.Client, pair string, interval string, limit int, smaPeriod int) bool {
	klines := GetKlines(client, pair, interval, limit)
	close := CloseFromKlines(klines)
	sma := SMA(close, smaPeriod)
	return IsMACrossOver(close, sma)
}

func RSI(prices []float64, period int) []float64 {
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

	// Étape 4 : Calculer le premier RS et le premier RSI
	// Le premier RSI est stocké à l'index `period`
	if avgLoss == 0 {
		rsi[period] = 100 // Pour éviter la division par zéro
	} else {
		rs := avgGain / avgLoss
		rsi[period] = 100 - (100 / (1 + rs))
	}

	// Étape 5 : Calculer les RSI suivants avec la méthode de lissage
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

func ExponetialMovingAverage(closingPrices []float64, period int) []float64 {

	if len(closingPrices) <= period {
		return SMA(closingPrices, period)
	}

	firstSMA := SMA(closingPrices, period)[0]
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
