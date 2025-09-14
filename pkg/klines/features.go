package klines

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Stupnikjs/binance_go/pkg/analysis"
	binance_connector "github.com/binance/binance-connector-go"
)

type Interval string

var Interv = []Interval{m1, m5, m15, m30, h1, h2, h4}

type Indicator struct {
	Name       string
	Interval   Interval
	Type       string
	Calculator func([]float64, int) []float64
	Param      int
}

func (i *Indicator) GetMapKey() string {
	return fmt.Sprintf("%s_%s_%d", i.Name, string(i.Interval), i.Param)
}

const (
	m1  Interval = "1m"
	m5  Interval = "5m"
	m15 Interval = "15m"
	m30 Interval = "30m"
	h1  Interval = "1h"
	h2  Interval = "2h"
	h4  Interval = "4h"
)

type FeaturedKlines struct {
	*binance_connector.KlinesResponse
	FeaturesMap map[string]float64
}

func GetEMAIndicatorsArray() [][]Indicator {
	var arr [][]Indicator

	for i := range 30 {

		for j := range 100 {
			if i < j && i != 0 && j != 0 {
				ar := []Indicator{
					{Name: "EMA", Interval: m15, Type: "Price", Calculator: analysis.EMAcalc, Param: i},
					{Name: "EMA", Interval: m15, Type: "Price", Calculator: analysis.EMAcalc, Param: j},
					{Name: "RSI", Interval: m15, Type: "Price", Calculator: analysis.RSIcalc, Param: 14},
				}
				arr = append(arr, ar)
			}
		}
	}
	return arr

}

func BuildIndicatorsMapArray(klines []*binance_connector.KlinesResponse, ind []Indicator) map[string][]float64 {
	mapper := make(map[string][]float64, len(ind))
	close := CloseFromKlines(klines)
	_ = VolumeFromKlines(klines)

	for _, i := range ind {

		if strings.HasPrefix(i.Name, "EMA") {
			mapper[i.GetMapKey()] = i.Calculator(close, i.Param)
		}

		if strings.HasPrefix(i.Name, "RSI") {
			mapper[i.GetMapKey()] = analysis.RSIcalc(close, i.Param)
		}
	}
	return mapper
}

func BuildFeaturedKlinesArray(klines []*binance_connector.KlinesResponse, ind []Indicator) []FeaturedKlines {
	var featuresArray []FeaturedKlines
	// THIS FUNC FAIL
	indicatorsArray := BuildIndicatorsMapArray(klines, ind)
	klen := len(klines)

	for i := 0; i < len(klines)-1; i++ {
		featured := FeaturedKlines{
			klines[i],
			map[string]float64{},
		}
		featured.FeaturesMap = make(map[string]float64, 1)

		for _, l := range ind {
			offset := klen - len(indicatorsArray[l.GetMapKey()])
			if i >= offset {
				featured.FeaturesMap[l.GetMapKey()] = indicatorsArray[l.GetMapKey()][i-offset]
			}
		}
		featuresArray = append(featuresArray, featured)

	}
	return featuresArray
}

func FeaturedKlinesToString(f FeaturedKlines) []string {
	var arr []string
	arr = append(arr, f.Close)
	arr = append(arr, f.Volume)
	strTime := strconv.Itoa(int(f.CloseTime))
	arr = append(arr, strTime)
	for _, v := range f.FeaturesMap {
		str := strconv.FormatFloat(v, 'f', 3, 64)
		arr = append(arr, str)
	}
	return arr
}

func (f *FeaturedKlines) FloatClose() float64 {

	float, err := strconv.ParseFloat(f.Close, 64)
	if err != nil {
		log.Fatal(err)
	}
	return float
}

func (f *FeaturedKlines) EMAShortOverLong(ind []Indicator) (bool, error) {
	if len(ind) < 2 {
		return false, fmt.Errorf("ind must be len 2 at least not %d", len(ind))
	}
	if !strings.HasPrefix(ind[0].Name, "EMA") || !strings.HasPrefix(ind[1].Name, "EMA") {
		return false, fmt.Errorf("two first indicators must be EMA not %s %s", ind[0].Name, ind[1].Name)
	}
	if f.FeaturesMap[ind[0].GetMapKey()] > f.FeaturesMap[ind[1].GetMapKey()] {
		return true, nil
	}
	return false, nil
}

func (f *FeaturedKlines) RSIUnder(threshold float64, ind []Indicator) bool {
	return f.FeaturesMap[ind[2].GetMapKey()] < threshold

}
