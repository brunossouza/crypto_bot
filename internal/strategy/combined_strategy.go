package strategy

import "github.com/brunossouza/crypto_bot/internal/indicators"

type CombinedStrategy struct {
	RSIPeriod          int
	SMAPeriod          int
	OverboughtLevel    float64
	OversoldLevel      float64
	TrendStrengthLevel float64
}

func NewCombinedStrategy(rsiPeriod, smaPeriod int, overbought, oversold, trendStrength float64) *CombinedStrategy {
	return &CombinedStrategy{
		RSIPeriod:          rsiPeriod,
		SMAPeriod:          smaPeriod,
		OverboughtLevel:    overbought,
		OversoldLevel:      oversold,
		TrendStrengthLevel: trendStrength,
	}
}

func (s *CombinedStrategy) ShouldEnter(prices []float64) bool {
	rsi := indicators.CalculateRSI(prices, s.RSIPeriod)
	sma := indicators.CalculateSMA(prices, s.SMAPeriod)
	currentPrice := prices[len(prices)-1]

	// Tendência de alta (preço acima da média móvel) + RSI indicando sobrevenda
	isTrendUp := currentPrice > sma
	isOversold := rsi < s.OversoldLevel
	trendStrength := (currentPrice - sma) / sma * 100

	return isTrendUp && isOversold && trendStrength > s.TrendStrengthLevel
}

func (s *CombinedStrategy) ShouldExit(prices []float64) bool {
	rsi := indicators.CalculateRSI(prices, s.RSIPeriod)
	sma := indicators.CalculateSMA(prices, s.SMAPeriod)
	currentPrice := prices[len(prices)-1]

	// Tendência de baixa (preço abaixo da média móvel) + RSI indicando sobrecompra
	isTrendDown := currentPrice < sma
	isOverbought := rsi > s.OverboughtLevel
	trendStrength := (sma - currentPrice) / sma * 100

	return isTrendDown && isOverbought && trendStrength > s.TrendStrengthLevel
}

func (s *CombinedStrategy) GetIndicators(prices []float64) (rsi, sma float64) {
	return indicators.CalculateRSI(prices, s.RSIPeriod),
		indicators.CalculateSMA(prices, s.SMAPeriod)
}
