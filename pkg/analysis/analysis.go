package analysis

import (
	"math"
)

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

// VROCcalc calculates the Volume Rate of Change indicator.
func VROCcalc(volumes []float64, period int) []float64 {

	// A period of 0 is invalid.
	if period <= 0 {
		return nil
	}
	// The number of volumes must be greater than the period.
	if len(volumes) <= period {
		return make([]float64, 0)
	}

	// as with moving average array is smaller
	vroc := make([]float64, len(volumes)-period)

	for i := period; i < len(volumes); i++ {
		currentVolume := volumes[i]
		previousVolume := volumes[i-period]

		// Avoid division by zero. If the previous volume is zero, the rate of change is 0.
		if previousVolume == 0 {
			vroc[i-period] = 0.0
			continue
		}

		// Calculate the VROC using the formula: (current - previous) / previous * 100
		vroc[i-period] = ((currentVolume - previousVolume) / previousVolume) * 100
	}

	return vroc
}

// VOLUME EMA
