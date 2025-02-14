package trading

import (
	"testing"
)

// TestParseFloat valida a conversão correta de string para float64.
func TestParseFloat(t *testing.T) {
	// ...existing code...
	val := parseFloat("1.23")
	if val != 1.23 {
		t.Errorf("Esperado 1.23, obteve %f", val)
	}
}

// TestCalculateRSI utiliza uma sequência de preços onde o RSI deve estar entre 0 e 100.
func TestCalculateRSI(t *testing.T) {
	// Cria uma série de preços simples
	prices := []float64{1.0, 2.0, 1.5, 2.5, 2.0, 3.0, 2.5, 3.5, 3.0, 4.0, 3.5, 4.5, 4.0, 5.0, 4.5}
	period := 14
	rsi := CalculateRSI(prices, period)
	if rsi < 0 || rsi > 100 {
		t.Errorf("RSI fora do intervalo esperado: %f", rsi)
	}
}

// TestCalculateRSIInsuficientPrices valida que um panic é disparado quando não há preços suficientes.
func TestCalculateRSIInsuficientPrices(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Era esperado panic por preços insuficientes")
		}
	}()
	// Utilize uma lista com exatamente 2 elementos para um período maior que 1.
	prices := []float64{1.0, 2.0}
	// Para período 5, é esperado panic, pois len(prices) < 6.
	_ = CalculateRSI(prices, 5)
}

// TestCalculateRSIPricesListSmallerThanPeriod valida que ocorre panic quando a quantidade de preços é menor que period + 1.
func TestCalculateRSIPricesListSmallerThanPeriod(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("Era esperado panic por lista de preços menor que o período requerido")
		}
	}()
	// Para período 14, uma lista com 10 elementos é insuficiente.
	prices := []float64{1.0, 1.2, 1.4, 1.3, 1.5, 1.6, 1.55, 1.7, 1.8, 1.75}
	_ = CalculateRSI(prices, 14)
}

// BenchmarkCalculateRSI mede o desempenho da função CalculateRSI.
func BenchmarkCalculateRSI(b *testing.B) {
	// Cria uma lista de preços com 1000 elementos
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
