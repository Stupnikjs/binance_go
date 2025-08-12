package main

import (
	"context"
	"fmt"
	"math"
	"strconv"

	binance_connector "github.com/binance/binance-connector-go"
)

const (
	RSI Indicator = "RSI"
	VOL Indicator = "VOL"
	SMA Indicator = "SMA"
	EMA Indicator = "EMA"
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

 var Interv = []Interval{m1,m5, m15, m30, h1, h2, h4 }
func findIntervalIndex(inter Interval) (int, ,error) {
	for i, in := range Interv {
		if in == inter {
			return i 
		}
	}
	return fmt.Errorf("Error occurs in findIntervalIndex")
}

type Indicators map[Indicator]float64

type MapKline map[Interval][]*binance_connector.KlinesResponse

type Kline struct {
	MainInterval  Interval
	Kline_binance *binance_connector.KlinesResponse
	IndicatorsMap map[Interval]Indicators
}

// Kline get upper interval
// query upper Intervals Coefs ex RSI_1h ..

// volume weighted Average Price

func BuildMapKline(client *binance_connector.Client, pair string) MapKline {
	mapK := make(MapKline)
	for _, i := range Interv {
		klines, err := client.NewKlinesService().
			Symbol(pair).
			Interval(string(i)).
			Limit(1000).
			Do(context.Background())
		if err != nil {
			fmt.Println(err)
		}
		mapK[i] = klines

	}
	return mapK
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

// Kline query upper period coef

// refactor
func IndicatorstoKlines(mapK MapKline, interval Interval, smallPeriod int, bigPeriod int, RSIcoef int) []Kline {
	klineArr := []Kline{}
	mainklines := mapK[interval]
	close := CloseFromKlines(mainklines)

	small_SMA := SMAcalc(close, smallPeriod)
	big_SMA := SMAcalc(close, bigPeriod)
	small_EMA := EMAcalc(close, smallPeriod)
	big_EMA := EMAcalc(close, bigPeriod)
	RSIcalc := RSIcalc(close, RSIcoef)

	small_period_index := -smallPeriod
	big_period_index := -bigPeriod
	for i, k := range mainklines {
		small_period_index += 1
		big_period_index += 1
		kl := Kline{
			Kline_binance: k,
			MainInterval:  interval,
			IndicatorsMap: make(map[Interval]Indicators),
		}
		// calulate the RSI or else of the klines up in timeframe that include this candle timeframe
		// recalculate only if close time of upper kline is bigger than  close time of k  
		upperIntervalIndex := findIntervalIndex(kl.MainInterval) + 1
		upperInterval := Interv[upperIntervalIndex]
		upperKlines := mapK[upperInterval] 
		
		RSIupper := RSIcalc(CloseFromKlines(upperKlines))
		

		kl.IndicatorsMap[kl.MainInterval]["RSI"] = RSIcalc[i]
		if small_period_index >= 0 {
			kl.IndicatorsMap[kl.MainInterval][SMA] = small_SMA[small_period_index]
			kl.IndicatorsMap[kl.MainInterval][EMA] = small_EMA[small_period_index]

			if big_period_index >= 0 {
				kl.IndicatorsMap[kl.MainInterval][SMA] = big_SMA[big_period_index]
				kl.IndicatorsMap[kl.MainInterval][EMA] = big_EMA[big_period_index]

			}

		}

		klineArr = append(klineArr, kl)

	}
	return klineArr
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
