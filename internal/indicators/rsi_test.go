package indicators

import "testing"

func TestCalculateRSI(t *testing.T) {
	prices := []float64{1.0, 2.0, 1.5, 2.5, 2.0, 3.0, 2.5, 3.5, 3.0, 4.0, 3.5, 4.5, 4.0, 5.0, 4.5}
	period := 14
	rsi := CalculateRSI(prices, period)
	if rsi < 0 || rsi > 100 {
		t.Errorf("RSI fora do intervalo esperado: %f", rsi)
	}
}

func TestCalculateRSIInsuficientPrices(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Era esperado panic por preços insuficientes")
		}
	}()
	prices := []float64{1.0, 2.0}
	_ = CalculateRSI(prices, 5)
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
