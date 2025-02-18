package utils

import (
	"testing"
)

// TestParseFloat valida a conversão correta de string para float64.
func TestParseFloat(t *testing.T) {
	// ...existing code...
	val := ParseFloat("1.23")
	if val != 1.23 {
		t.Errorf("Esperado 1.23, obteve %f", val)
	}
}
