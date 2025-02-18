package indicators

import (
	"testing"
)

func TestCalculateSMA(t *testing.T) {
	tests := []struct {
		name      string
		prices    []float64
		period    int
		expected  float64
		wantPanic bool
	}{
		{
			name:      "Should calculate SMA correctly for period 3",
			prices:    []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			period:    3,
			expected:  4.0, // (3 + 4 + 5) / 3
			wantPanic: false,
		},
		{
			name:      "Should calculate SMA correctly for period 5",
			prices:    []float64{1.0, 2.0, 3.0, 4.0, 5.0},
			period:    5,
			expected:  3.0, // (1 + 2 + 3 + 4 + 5) / 5
			wantPanic: false,
		},
		{
			name:      "Should panic when prices length is less than period",
			prices:    []float64{1.0, 2.0},
			period:    3,
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

			got := CalculateSMA(tt.prices, tt.period)

			if !tt.wantPanic && got != tt.expected {
				t.Errorf("CalculateSMA() = %v, want %v", got, tt.expected)
			}
		})
	}
}
