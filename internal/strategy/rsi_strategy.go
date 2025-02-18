package strategy

import "github.com/brunossouza/crypto_bot/internal/indicators"

type RSIStrategy struct {
	Period          int
	OverboughtLevel float64
	OversoldLevel   float64
}

func NewRSIStrategy(period int, overbought, oversold float64) *RSIStrategy {
	return &RSIStrategy{
		Period:          period,
		OverboughtLevel: overbought,
		OversoldLevel:   oversold,
	}
}

func (s *RSIStrategy) ShouldEnter(prices []float64) bool {
	rsi := indicators.CalculateRSI(prices, s.Period)
	return rsi < s.OversoldLevel
}

func (s *RSIStrategy) ShouldExit(prices []float64) bool {
	rsi := indicators.CalculateRSI(prices, s.Period)
	return rsi > s.OverboughtLevel
}

func (s *RSIStrategy) GetRSI(prices []float64) float64 {
	return indicators.CalculateRSI(prices, s.Period)
}
