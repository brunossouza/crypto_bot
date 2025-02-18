package indicators

// CalculateSMA calcula a Média Móvel Simples para uma série de preços
// Parâmetros:
//   - prices: slice com os preços históricos
//   - period: período para cálculo da média
//
// Retorna:
//   - float64: valor da média móvel
func CalculateSMA(prices []float64, period int) float64 {
	if len(prices) < period {
		panic("Not enough prices to calculate SMA")
	}

	sum := 0.0
	for i := len(prices) - period; i < len(prices); i++ {
		sum += prices[i]
	}

	return sum / float64(period)
}
