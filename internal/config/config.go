package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/joho/godotenv"
)

// Config armazena as configurações necessárias para o funcionamento do bot de criptomoedas.
// Contém informações de conexão com a API da Binance e parâmetros de negociação.
type Config struct {
	// ApiURL é o endpoint base da API da Binance
	ApiURL string
	// Symbol é o par de criptomoedas para negociação (ex: BTCUSDT)
	Symbol string
	// Period é o intervalo em minutos para análise do mercado
	Period int
	// ApiKey é a chave pública da API da Binance
	ApiKey string
	// ApiSecret é a chave privada da API da Binance
	ApiSecret string
}

// LoadConfig carrega as configurações do arquivo .env e valida os valores obrigatórios.
// O método verifica:
// - Se o arquivo .env pode ser carregado
// - Se as credenciais da API (BINANCE_API_KEY e BINANCE_API_SECRET) estão presentes
// - Se os parâmetros básicos (API_URL, SYMBOL e PERIOD) estão configurados corretamente
// - Se o valor de PERIOD é um número inteiro válido
//
// Retorna:
// - Um ponteiro para Config com as configurações carregadas
// - Um erro se houver falha no carregamento ou validação
func LoadConfig() (*Config, error) {
	// Tenta carregar o arquivo .env, mas não interrompe a execução se não existir
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found, using environment variables")
	}

	fmt.Println("Carregando configurações...")

	conf := &Config{
		ApiURL:    os.Getenv("API_URL"),
		Symbol:    os.Getenv("SYMBOL"),
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET"),
	}

	var missingVars []string
	if conf.ApiURL == "" {
		missingVars = append(missingVars, "API_URL")
	}
	if conf.Symbol == "" {
		missingVars = append(missingVars, "SYMBOL")
	}
	if conf.ApiKey == "" {
		missingVars = append(missingVars, "BINANCE_API_KEY")
	}
	if conf.ApiSecret == "" {
		missingVars = append(missingVars, "BINANCE_API_SECRET")
	}

	periodStr := os.Getenv("PERIOD")
	period, err := strconv.Atoi(periodStr)
	if err != nil {
		missingVars = append(missingVars, "PERIOD")
	}
	conf.Period = period

	if len(missingVars) > 0 {
		return nil, fmt.Errorf("as variáveis obrigatórias estão faltando: %s", strings.Join(missingVars, ", "))
	}

	fmt.Println("Configurações carregadas com sucesso")

	return conf, nil
}
