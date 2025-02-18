package utils

import (
	"log"
	"strconv"
)

// parseFloat converte uma string para float64
// Parâmetros:
// - str: string contendo um número decimal
//
// O método:
// 1. Tenta converter a string para float64
// 2. Em caso de erro, finaliza o programa
//
// Retorna:
// - float64: valor numérico convertido da string
func ParseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return val
}
