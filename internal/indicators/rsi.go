package indicators

import "math"

// CalculateRSI calcula o Índice de Força Relativa (RSI) para uma série de preços
// O RSI é um indicador de momentum que mede a velocidade e magnitude das mudanças de preços
// para identificar condições de sobrecompra ou sobrevenda.
//
// Parâmetros:
//   - prices: slice com os preços históricos ordenados do mais antigo para o mais recente
//   - period: período para o cálculo do RSI (geralmente 14 períodos)
//
// Retorna:
//   - float64: valor do RSI entre 0 e 100
//   - Valores acima de 70 geralmente indicam sobrecompra
//   - Valores abaixo de 30 geralmente indicam sobrevenda
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		panic("Not enough prices to calculate RSI")
	}

	var avgGains, avgLoss float64

	// Calcula valores iniciais
	for i := 1; i < len(prices); i++ {
		gain, loss := calculateAverage(prices, period, i)

		if i == 1 {
			avgGains = gain
			avgLoss = loss
			continue
		}

		// Aplica a fórmula da Média Móvel Exponencial (EMA):
		// EMA = (Valor atual * (1/N)) + (EMA anterior * (N-1)/N)
		// Onde N é o período e EMA é calculado separadamente para ganhos e perdas
		avgGains = (avgGains*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
	}

	rs := avgGains / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

// calculateAverage calcula a média de ganhos e perdas para um determinado período
// Esta função é utilizada internamente pelo CalculateRSI para processar os dados
//
// Parâmetros:
//   - prices: slice com os preços históricos
//   - period: período para cálculo da média
//   - startIdx: índice inicial para começar o cálculo
//
// Retorna:
//   - float64: média dos ganhos no período
//   - float64: média das perdas no período
func calculateAverage(prices []float64, period int, startIdx int) (float64, float64) {
	var gain, loss float64
	for i := 0; i < period && i+startIdx < len(prices); i++ {
		diff := prices[i+startIdx] - prices[i+startIdx-1]
		if diff > 0 {
			gain += diff
		} else {
			loss += math.Abs(diff)
		}
	}
	return gain / float64(period), loss / float64(period)
}
