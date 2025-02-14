package main

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

// Config armazena todas as configurações do bot
type Config struct {
	ApiURL    string
	Symbol    string
	Period    int
	ApiKey    string
	ApiSecret string
}

// LoadConfig carrega e valida todas as configurações do arquivo .env
func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("erro ao carregar arquivo .env: %v", err)
	}

	config := &Config{
		ApiURL:    os.Getenv("API_URL"),
		Symbol:    os.Getenv("SYMBOL"),
		ApiKey:    os.Getenv("BINANCE_API_KEY"),
		ApiSecret: os.Getenv("BINANCE_API_SECRET"),
	}

	// Converte PERIOD para inteiro
	periodStr := os.Getenv("PERIOD")
	period, err := strconv.Atoi(periodStr)
	if err != nil {
		return nil, fmt.Errorf("PERIOD deve ser um número inteiro válido: %v", err)
	}
	config.Period = period

	// Validação das configurações
	if config.ApiKey == "" || config.ApiSecret == "" {
		return nil, fmt.Errorf("as variáveis de ambiente BINANCE_API_KEY e BINANCE_API_SECRET são obrigatórias")
	}
	if config.ApiURL == "" || config.Symbol == "" || config.Period <= 0 {
		return nil, fmt.Errorf("as variáveis de ambiente API_URL, SYMBOL e PERIOD são obrigatórias")
	}

	return config, nil
}

// Atualiza as variáveis globais para usar config
var (
	IsOpened bool = false
	config   *Config
)

type Candlestick struct {
	OpenTime                 int64   // Horário de abertura do candle em milissegundos
	Open                     float64 // Preço de abertura do candle
	High                     float64 // Preço mais alto durante o período do candle
	Low                      float64 // Preço mais baixo durante o período do candle
	Close                    float64 // Preço de fechamento (ou último preço) do candle
	Volume                   float64 // Volume total negociado durante o período do candle
	CloseTime                int64   // Horário de fechamento do candle em milissegundos
	QuoteAssetVolume         float64 // Volume total do ativo de cotação durante o período
	NumberOfTrades           int64   // Número de negociações durante o período
	TakerBuyBaseAssetVolume  float64 // Volume de compra do ativo base pelos takers
	TakerBuyQuoteAssetVolume float64 // Volume de compra do ativo de cotação pelos takers
	Ignore                   float64 // Campo não utilizado, ignorar
}

// Update GetCandlesticks to use config variable
func GetCandlesticks(symbol string, interval string, limit int) []Candlestick {
	// Cria uma nova requisição
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", config.ApiURL, symbol, interval, limit), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Lê o corpo da resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Processa a resposta JSON
	var rawData [][]interface{}
	err = json.Unmarshal(body, &rawData)
	if err != nil {
		log.Fatal(err)
	}

	candlesticks := make([]Candlestick, len(rawData))
	for i, raw := range rawData {
		candlesticks[i] = Candlestick{
			OpenTime:                 int64(raw[0].(float64)),
			Open:                     parseFloat(raw[1].(string)),
			High:                     parseFloat(raw[2].(string)),
			Low:                      parseFloat(raw[3].(string)),
			Close:                    parseFloat(raw[4].(string)),
			Volume:                   parseFloat(raw[5].(string)),
			CloseTime:                int64(raw[6].(float64)),
			QuoteAssetVolume:         parseFloat(raw[7].(string)),
			NumberOfTrades:           int64(raw[8].(float64)),
			TakerBuyBaseAssetVolume:  parseFloat(raw[9].(string)),
			TakerBuyQuoteAssetVolume: parseFloat(raw[10].(string)),
			Ignore:                   parseFloat(raw[11].(string)),
		}
	}

	return candlesticks
}

// parseFloat converte uma string para float64
// Função auxiliar utilizada para converter os valores string da API para números
// Em caso de erro na conversão, finaliza o programa
func parseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// calculateAverage calcula a média de ganhos e perdas para um determinado período
// prices: slice com os preços históricos
// period: período para cálculo da média
// startIdx: índice inicial para começar o cálculo
// Retorna dois float64: média de ganhos e média de perdas do período
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

// CalculateRSI calcula o Índice de Força Relativa (RSI) usando Média Móvel Exponencial (EMA)
// prices: slice com os preços históricos
// period: período para cálculo do RSI (geralmente 14)
// Utiliza EMA para dar mais peso aos preços recentes
// Retorna o valor do RSI entre 0 e 100
// - Valores acima de 70 indicam sobrecompra (overbought)
// - Valores abaixo de 30 indicam sobrevenda (oversold)
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		log.Fatal("Not enough prices to calculate RSI")
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

		// Calcula EMA (Exponential Moving Average) para ganhos e perdas
		avgGains = (avgGains*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
	}

	rs := avgGains / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

// Update NewOrder to use config variable
func NewOrder(symbol string, quantity float64, side string) error {
	// Prepara os parâmetros da ordem
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("quantity", fmt.Sprintf("%f", quantity))
	params.Add("side", side)
	params.Add("type", "MARKET")
	params.Add("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// Gera a assinatura HMAC SHA256
	mac := hmac.New(sha256.New, []byte(config.ApiSecret))
	mac.Write([]byte(params.Encode()))
	signature := hex.EncodeToString(mac.Sum(nil))
	params.Add("signature", signature)

	// Cria a requisição
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v3/order", config.ApiURL), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	// Adiciona os headers necessários
	req.Header.Add("X-MBX-APIKEY", config.ApiKey)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	// Envia a requisição
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Lê a resposta
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Se o status não for 200, retorna erro
	if resp.StatusCode != 200 {
		return fmt.Errorf("erro na criação da ordem: %s", string(body))
	}

	fmt.Printf("Ordem criada com sucesso: %s\n", string(body))
	return nil
}

// Update StartTrading to use config variable
func StartTrading() {
	// Obtém os dados dos candles
	candlesticks := GetCandlesticks(config.Symbol, "15m", 100)

	// Obtém o último preço
	lastPrice := candlesticks[len(candlesticks)-1].Close

	var prices []float64
	for _, candlestick := range candlesticks {
		prices = append(prices, candlestick.Close)
	}

	// Calcula o RSI
	rsi := CalculateRSI(prices, config.Period)

	// Limpa a tela
	fmt.Print("\033[H\033[2J")

	// Imprime o último preço
	fmt.Printf("Último preço: %.2f\n", lastPrice)
	// Imprime o RSI
	fmt.Printf("RSI: %.2f\n", rsi)

	if rsi < 30 && !IsOpened {
		fmt.Println("sobrevendido, momento de comprar")
		if err := NewOrder(config.Symbol, 0.001, "BUY"); err != nil {
			log.Println(err)
			IsOpened = false
		} else {
			IsOpened = true
		}
	} else if rsi > 70 && IsOpened {
		fmt.Println("sobrecomprado, momento de vender")
		if err := NewOrder(config.Symbol, 0.001, "SELL"); err != nil {
			log.Println(err)
			IsOpened = true
		} else {
			IsOpened = false
		}
	} else {
		fmt.Println("Aguardando oportunidades...")
	}
}

// main é o ponto de entrada do programa
// Configura um temporizador para executar StartTrading a cada 3 segundos
// Continua executando até o programa ser interrompido (CTRL+C)
func main() {
	var err error
	config, err = LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Cria um temporizador que dispara a cada 10 segundos
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Bot iniciado! Pressione CTRL+C para parar")

	StartTrading()

	// Executa indefinidamente até o programa ser interrompido
	for range ticker.C {
		StartTrading()
	}
}
