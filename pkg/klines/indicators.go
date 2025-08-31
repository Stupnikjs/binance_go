package klines

import (
	"fmt"

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

var Indicators = []Indicator{
	{"RSI", m5, "Price", analysis.RSIcalc, 14},
	{"EMA", m5, "Price", analysis.EMAcalc, 14},
}

func BuildIndicatorsMapArray(klines []*binance_connector.KlinesResponse, ind []Indicator) map[string][]float64 {
	mapper := make(map[string][]float64, len(Indicators))
	close := CloseFromKlines(klines)
	vols := VolumeFromKlines(klines)

	for _, i := range Indicators {
		if i.Type == "Price" {
			mapper[i.GetMapKey()] = i.Calculator(close, i.Param)
		}
		if i.Type == "Volume" {
			mapper[i.GetMapKey()] = i.Calculator(vols, i.Param)
		}
	}
	return mapper
}

func BuildFeaturedKlinesArray(klines []*binance_connector.KlinesResponse, ind []Indicator) []FeaturedKlines {
	var featuresArray []FeaturedKlines
	indicatorsArray := BuildIndicatorsMapArray(klines, ind)
	klen := len(klines)
	for i, k := range klines {
		featured := FeaturedKlines{
			k,
			make(map[string]float64),
		}
		for _, l := range ind {
			if len(indicatorsArray[l.GetMapKey()])-klen+i > 0 {
				featured.FeaturesMap[l.GetMapKey()] = indicatorsArray[l.GetMapKey()][i]
			}
		}
		featuresArray = append(featuresArray, featured)

	}
	return featuresArray
}
