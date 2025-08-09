package main

import "testing"

func TestSimpleMovingAverage(t *testing.T) {
	type testCase struct {
		name        string
		prices      []float64
		period      int
		expectedSMA []float64
	}

	tests := []testCase{
		{
			name:        "Basic 2-period SMA",
			prices:      []float64{20, 15, 22, 30, 40, 50, 60},
			period:      2,
			expectedSMA: []float64{17.5, 18.5, 26, 35, 45, 55}, // Corrected expected values
		},
		{
			name:        "Basic 3-period SMA",
			prices:      []float64{10, 20, 30, 40, 50, 60, 70, 80},
			period:      3,
			expectedSMA: []float64{20, 30, 40, 50, 60, 70},
		},
		{
			name:        "Period larger than prices",
			prices:      []float64{10, 20},
			period:      3,
			expectedSMA: []float64{},
		},
		{
			name:        "Empty price slice",
			prices:      []float64{},
			period:      2,
			expectedSMA: []float64{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actualSMA := SMA(tt.prices, tt.period)
			if len(actualSMA) != len(tt.expectedSMA) {
				t.Fatalf("lengths do not match: expected %v, got %v", tt.expectedSMA, actualSMA)
			}
			for i := range actualSMA {
				if actualSMA[i] != tt.expectedSMA[i] {
					t.Errorf("at index %d, expected %v, got %v", i, tt.expectedSMA[i], actualSMA[i])
				}
			}
		})
	}
}
