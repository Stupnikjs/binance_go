package klines

import (
	"context"
	"fmt"
	"path"
	"strconv"
	"strings"
	"time"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	binance_connector "github.com/binance/binance-connector-go"
)

type Indicator string

const (
	RSI     Indicator = "RSI"
	VOL     Indicator = "VOL"
	VROC    Indicator = "VROC"
	RSI_15m Indicator = "RSI_15m"
	RSI_30m Indicator = "RSI_30m"
	RSI_1h  Indicator = "RSI_1h"
	RSI_2h  Indicator = "RSI_2h"
	RSI_4h  Indicator = "RSI_4h"

	SMA_short      Indicator = "SMA_short"
	EMA_short      Indicator = "EMA_short"
	SMA_long       Indicator = "SMA_long"
	EMA_long       Indicator = "EMA_long"
	SMA_super_long Indicator = "SMA_super_long"
	EMA_super_long Indicator = "EMA_super_long"
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

func IntervalToTime(interval Interval) (time.Duration, error) {
	// The time.ParseDuration function can handle strings like "10m", "1h", "1m30s".
	// It returns a duration and an error.
	return time.ParseDuration(string(interval))
}

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type Indicators map[Indicator][]float64
type Klines struct {
	Interval   Interval
	Array      []*binance_connector.KlinesResponse
	Indicators Indicators
}

type IndicatorsParams struct {
	Short_period_MA int
	Long_period_MA  int
	Super_long_MA   int
	RSI_coef        int
	VROC_coef       int
}

func BuildKlineArrData(pair string, interval []Interval) []*Klines {
	klinesArr := []*Klines{}
	binanceArr := []*binance_connector.KlinesResponse{}
	path := path.Join("data", string(interval[0]), strings.ToLower(pair))
	kline, err := LoadKlinesFromFile(path)
	if err != nil {
		fmt.Println(err)
	}
	for _, k := range kline {
		binanceArr = append(binanceArr, &k)
	}
	kl := Klines{
		Array:      binanceArr,
		Interval:   interval[0],
		Indicators: make(Indicators),
	}

	klinesArr = append(klinesArr, &kl)
	return klinesArr
}

func BuildKlinesArrLive(client *binance_connector.Client, pair string, Interval []Interval) []*Klines {
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

func VolumeFromKlines(klines []*binance_connector.KlinesResponse) []float64 {
	volumes := make([]float64, len(klines))
	for i, kline := range klines {
		f_vol, err := strconv.ParseFloat(kline.Volume, 64)
		if err != nil {
			fmt.Println(err)
		}
		volumes[i] = f_vol

	}
	return volumes
}

// error checking
func IndicatorstoKlines(client *binance_connector.Client, pair string, intervals []Interval, params IndicatorsParams) []*Klines {
	klinesArr := BuildKlinesArrLive(client, pair, intervals)
	return ProcessKlinesNormalizedRefactored(klinesArr, params)
}

func ProcessKlinesNormalized(klines []*Klines, params IndicatorsParams) []*Klines {
	for _, k := range klines {
		// caclulate RSI EMA SMA
		close := CloseFromKlines(k.Array)
		vols := VolumeFromKlines(k.Array)
		RSI_arr := analysis.RSIcalc(close, params.RSI_coef)
		VROC_arr := analysis.VROCcalc(vols, params.VROC_coef)
		SMA_short_arr := analysis.SMAcalc(close, params.Short_period_MA)
		SMA_long_arr := analysis.SMAcalc(close, params.Long_period_MA)
		EMA_short_arr := analysis.EMAcalc(close, params.Short_period_MA)
		EMA_long_arr := analysis.EMAcalc(close, params.Long_period_MA)
		SMA_super_long_arr := analysis.SMAcalc(close, params.Super_long_MA)
		EMA_super_long_arr := analysis.EMAcalc(close, params.Super_long_MA)

		// return sliced array of same length
		offset := params.Super_long_MA
		k.Array = k.Array[offset-1:]
		k.Indicators[RSI] = RSI_arr[offset:]
		k.Indicators[VROC] = VROC_arr[offset-params.VROC_coef-1:]
		k.Indicators[SMA_short] = SMA_short_arr[offset-params.Short_period_MA:]
		k.Indicators[SMA_long] = SMA_long_arr[offset-params.Long_period_MA:]
		k.Indicators[EMA_short] = EMA_short_arr[offset-params.Short_period_MA:]
		k.Indicators[EMA_long] = EMA_long_arr[offset-params.Long_period_MA:]
		k.Indicators[SMA_super_long] = SMA_super_long_arr
		k.Indicators[EMA_super_long] = EMA_super_long_arr

	}

	// MeltRSIKline(klines[0], klines[2]) bugging with vroc
	return klines
}

// origin must be klines from upper interval
func MeltRSIKline(receiver *Klines, origin *Klines) error {
	originRSI := origin.Indicators[RSI]
	if len(originRSI) == 0 {
		return fmt.Errorf("origin RSI indicator is empty or not calculated")
	}

	var targetIndicator Indicator
	// switch on Intervals
	switch origin.Interval {
	case m15:
		targetIndicator = RSI_15m
	case m30:
		targetIndicator = RSI_30m
	case h1:
		targetIndicator = RSI_1h
	case h2:
		targetIndicator = RSI_2h
	case h4:
		targetIndicator = RSI_4h
	default:
		return fmt.Errorf("no indicator valid found")
	}

	// Initialize the slice to the correct size
	receiver.Indicators[targetIndicator] = make([]float64, len(receiver.Array))

	originIndex := 0
	// Find the first corresponding origin candle
	for i := range origin.Array {
		if receiver.Array[0].OpenTime >= origin.Array[i].OpenTime {
			originIndex = i
		}
	}

	// Iterate over the receiver array (the lower interval)
	for i := range receiver.Array {
		// Advance the originIndex to find the correct high-interval candle.
		for originIndex+1 < len(origin.Array) && receiver.Array[i].CloseTime >= origin.Array[originIndex+1].CloseTime {
			originIndex++
		}

		// Assign the RSI value from the correctly identified origin candle
		receiver.Indicators[targetIndicator][i] = originRSI[originIndex]
	}

	return nil
}

func ProcessKlinesNormalizedRefactored(klines []*Klines, params IndicatorsParams) []*Klines {
	// Longest period determines the final length of all arrays.
	maxPeriod := params.Super_long_MA

	for _, k := range klines {
		close_arr := CloseFromKlines(k.Array)
		vols_arr := VolumeFromKlines(k.Array)

		// Define indicator calculations and their parameters in a structured way.
		indicatorsToCalc := map[Indicator]struct {
			calcFunc func([]float64, int) []float64
			data     []float64
			period   int
		}{
			RSI:            {analysis.RSIcalc, close_arr, params.RSI_coef},
			VROC:           {analysis.VROCcalc, vols_arr, params.VROC_coef - 1},
			SMA_short:      {analysis.SMAcalc, close_arr, params.Short_period_MA},
			SMA_long:       {analysis.SMAcalc, close_arr, params.Long_period_MA},
			EMA_short:      {analysis.EMAcalc, close_arr, params.Short_period_MA},
			EMA_long:       {analysis.EMAcalc, close_arr, params.Long_period_MA},
			SMA_super_long: {analysis.SMAcalc, close_arr, params.Super_long_MA},
			EMA_super_long: {analysis.EMAcalc, close_arr, params.Super_long_MA},
		}
		// Slice the original Klines array first.
		k.Array = k.Array[maxPeriod:]

		// Loop through the defined calculations, perform them, and store the sliced results.
		for name, ind := range indicatorsToCalc {
			result := ind.calcFunc(ind.data, ind.period)
			k.Indicators[name] = result[maxPeriod-ind.period:]
		}
	}

	return klines
}

func (k *Klines) SMAShortOverLong(index int) bool {
	return k.Indicators[SMA_short][index] < k.Indicators[SMA_long][index]
}
func (k *Klines) EMAShortOverLong(index int) bool {
	return k.Indicators[EMA_short][index] < k.Indicators[EMA_long][index]
}

func ConvertUSDCtoPAIR(client *binance_connector.Client, USDCamount float64, pair string) float64 {
	klines := BuildKlinesArr(client, pair, []Interval{m1})
	f_close, err := strconv.ParseFloat(klines[0].Array[0].Close, 64)
	if err != nil {
		fmt.Println(err)
	}
	return USDCamount / f_close
}
