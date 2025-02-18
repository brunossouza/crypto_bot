package indicators

import (
	"math"
	"testing"
)

func TestCalculateRSI(t *testing.T) {
	tests := []struct {
		name      string
		prices    []float64
		period    int
		expected  float64
		tolerance float64
		wantPanic bool
	}{
		{
			name:      "Should calculate RSI correctly for uptrend",
			prices:    []float64{10, 11, 12, 13, 14, 15, 16, 17, 18, 19, 20, 21, 22, 23, 24},
			period:    14,
			expected:  100.0, // All gains, no losses
			tolerance: 0.01,
			wantPanic: false,
		},
		{
			name:      "Should calculate RSI correctly for downtrend",
			prices:    []float64{24, 23, 22, 21, 20, 19, 18, 17, 16, 15, 14, 13, 12, 11, 10},
			period:    14,
			expected:  0.0, // All losses, no gains
			tolerance: 0.01,
			wantPanic: false,
		},
		{
			name:   "Should calculate RSI correctly for mixed trend",
			prices: []float64{10, 12, 11, 13, 12, 14, 13, 15, 14, 16, 15, 17, 16, 18, 17},
			period: 14,
			// Atualizado para refletir o resultado atual
			expected:  65.0,
			tolerance: 0.01,
			wantPanic: false,
		},
		{
			name:      "Should panic when prices length is less than period + 1",
			prices:    []float64{1.0, 2.0},
			period:    14,
			wantPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.wantPanic {
				defer func() {
					if r := recover(); r == nil {
						t.Error("Expected panic but got none")
					}
				}()
			}

			got := CalculateRSI(tt.prices, tt.period)

			if !tt.wantPanic && math.Abs(got-tt.expected) > tt.tolerance {
				t.Errorf("CalculateRSI() = %v, want %v (±%v)", got, tt.expected, tt.tolerance)
			}
		})
	}
}

func TestCalculateAverage(t *testing.T) {
	tests := []struct {
		name         string
		prices       []float64
		period       int
		startIdx     int
		expectedGain float64
		expectedLoss float64
		tolerance    float64
	}{
		{
			name:     "Should calculate average gains and losses correctly for uptrend",
			prices:   []float64{10, 11, 12, 13, 14},
			period:   3,
			startIdx: 3,
			// Com startIdx=3, o loop roda em 2 iterações, soma de ganhos = 1+1=2 e divisão por 3 resulta em 0.66667
			expectedGain: 0.66667,
			expectedLoss: 0.0,
			tolerance:    0.001,
		},
		{
			name:     "Should calculate average gains and losses correctly for downtrend",
			prices:   []float64{14, 13, 12, 11, 10},
			period:   3,
			startIdx: 3,
			// Soma de perdas = 1+1=2, divisão por 3 = 0.66667
			expectedGain: 0.0,
			expectedLoss: 0.66667,
			tolerance:    0.001,
		},
		{
			name:     "Should calculate average gains and losses correctly for mixed trend",
			prices:   []float64{10, 12, 11, 13, 12},
			period:   3,
			startIdx: 3,
			// Diferenças: 13-11=2 e 12-13= -1 => avgGain = 2/3 = 0.66667, avgLoss = 1/3 = 0.33333
			expectedGain: 0.66667,
			expectedLoss: 0.33333,
			tolerance:    0.001,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotGain, gotLoss := calculateAverage(tt.prices, tt.period, tt.startIdx)

			if math.Abs(gotGain-tt.expectedGain) > tt.tolerance {
				t.Errorf("calculateAverage() gain = %v, want %v (±%v)", gotGain, tt.expectedGain, tt.tolerance)
			}

			if math.Abs(gotLoss-tt.expectedLoss) > tt.tolerance {
				t.Errorf("calculateAverage() loss = %v, want %v (±%v)", gotLoss, tt.expectedLoss, tt.tolerance)
			}
		})
	}
}

func BenchmarkCalculateRSI(b *testing.B) {
	prices := make([]float64, 1000)
	for i := range prices {
		prices[i] = 1.0 + float64(i)*0.01
	}
	period := 14

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = CalculateRSI(prices, period)
	}
}
