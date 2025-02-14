package trading

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
	"strconv"
	"strings"
	"time"

	"github.com/brunossouza/crypto_bot/internal/config"
	"github.com/brunossouza/crypto_bot/internal/database"
)

var (
	cfg      *config.Config
	IsOpened bool = false
)

// Initialize define as configurações para o pacote de trading
// Parâmetros:
// - c: ponteiro para a estrutura de configuração contendo as credenciais da API e parâmetros do bot
// O método armazena a configuração em uma variável global para uso em todo o pacote
func Initialize(c *config.Config) {
	cfg = c
	if err := database.Initialize(); err != nil {
		log.Fatal("Erro ao inicializar banco de dados:", err)
	}
	// Load last status from the database
	status, err := database.GetPosition(cfg.Symbol)
	if err != nil {
		log.Println("Erro ao carregar status da posição:", err)
	}
	IsOpened = status
}

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

// GetCandlesticks obtém os dados históricos de preços do par de moedas especificado
// Parâmetros:
// - symbol: par de moedas para obter dados (ex: BTCUSDT)
// - interval: intervalo de tempo entre cada candle (ex: "15m", "1h", "4h")
// - limit: quantidade máxima de candles a serem retornados
//
// O método:
// 1. Faz uma requisição GET para a API da Binance
// 2. Processa a resposta JSON
// 3. Converte os dados para a estrutura Candlestick
//
// Retorna:
// - []Candlestick: slice contendo os dados históricos formatados
func GetCandlesticks(symbol string, interval string, limit int) []Candlestick {
	// Cria uma nova requisição
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", cfg.ApiURL, symbol, interval, limit), nil)
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
// Parâmetros:
// - str: string contendo um número decimal
//
// O método:
// 1. Tenta converter a string para float64
// 2. Em caso de erro, finaliza o programa
//
// Retorna:
// - float64: valor numérico convertido da string
func parseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// calculateAverage calcula a média de ganhos e perdas para um determinado período
// Parâmetros:
// - prices: slice com os preços históricos
// - period: período para cálculo da média
// - startIdx: índice inicial para começar o cálculo
//
// O método:
// 1. Itera sobre os preços no período especificado
// 2. Calcula a diferença entre preços consecutivos
// 3. Acumula ganhos e perdas separadamente
//
// Retorna:
// - float64: média dos ganhos no período
// - float64: média das perdas no período
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

// CalculateRSI calcula o Índice de Força Relativa (RSI)
// Parâmetros:
// - prices: slice com os preços históricos
// - period: período para cálculo do RSI (geralmente 14)
//
// O método:
//  1. Verifica se há dados suficientes para o cálculo
//  2. Calcula as médias iniciais de ganhos e perdas
//  3. Aplica o cálculo da Média Móvel Exponencial (EMA)
//  4. Calcula o RSI usando a fórmula: 100 - (100 / (1 + RS))
//     onde RS = EMA dos ganhos / EMA das perdas
//
// Retorna:
// - float64: valor do RSI entre 0 e 100
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

		// Calcula EMA (Exponential Moving Average) para ganhos e perdas
		avgGains = (avgGains*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
	}

	rs := avgGains / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

// NewOrder cria uma nova ordem de compra ou venda no mercado
// Parâmetros:
// - symbol: par de moedas para negociação (ex: BTCUSDT)
// - quantity: quantidade do ativo a ser negociada
// - side: direção da ordem ("BUY" para compra, "SELL" para venda)
// - price: preço atual do ativo no momento da ordem
//
// Retorna:
// - error: nil em caso de sucesso, ou erro em caso de falha
func NewOrder(symbol string, quantity float64, side string, price float64) error {
	// Prepara os parâmetros da ordem
	params := url.Values{}
	params.Add("symbol", symbol)
	params.Add("quantity", fmt.Sprintf("%f", quantity))
	params.Add("side", side)
	params.Add("type", "MARKET")
	params.Add("timestamp", fmt.Sprintf("%d", time.Now().UnixMilli()))

	// Gera a assinatura HMAC SHA256
	mac := hmac.New(sha256.New, []byte(cfg.ApiSecret))
	mac.Write([]byte(params.Encode()))
	signature := hex.EncodeToString(mac.Sum(nil))
	params.Add("signature", signature)

	// Cria a requisição
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/api/v3/order", cfg.ApiURL), strings.NewReader(params.Encode()))
	if err != nil {
		return err
	}

	// Adiciona os headers necessários
	req.Header.Add("X-MBX-APIKEY", cfg.ApiKey)
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

	// Se a ordem foi criada com sucesso, salva no banco
	if err := database.SaveOrder(symbol, side, quantity, price); err != nil {
		return fmt.Errorf("erro ao salvar ordem: %v", err)
	}

	// Atualiza a posição no banco
	isOpened := side == "BUY"
	if err := database.UpdatePosition(symbol, isOpened); err != nil {
		return fmt.Errorf("erro ao atualizar posição: %v", err)
	}

	fmt.Printf("Ordem criada com sucesso: %s\n", string(body))
	return nil
}

// StartTrading executa a lógica principal de trading do bot
// O método:
// 1. Obtém os dados mais recentes dos candles
// 2. Extrai o último preço e histórico de preços
// 3. Calcula o RSI atual
// 4. Atualiza a interface com informações do mercado
// 5. Executa a lógica de trading baseada no RSI:
//   - Compra quando RSI < 30 (sobrevenda)
//   - Vende quando RSI > 70 (sobrecompra)
//
// Comportamento:
// - Mantém controle do estado da posição através da variável IsOpened
// - Executa ordens de mercado com quantidade fixa de 0.001
// - Exibe mensagens de status no console
func StartTrading() {
	// Obtém os dados dos candles
	candlesticks := GetCandlesticks(cfg.Symbol, "15m", 100)

	// Obtém o último preço
	lastPrice := candlesticks[len(candlesticks)-1].Close

	var prices []float64
	for _, c := range candlesticks {
		prices = append(prices, c.Close)
	}

	// Calcula o RSI
	rsi := CalculateRSI(prices, cfg.Period)

	// Limpa a tela
	fmt.Print("\033[H\033[2J")
	fmt.Println("API URL:", cfg.ApiURL)
	fmt.Println("Ativo:", cfg.Symbol)
	fmt.Printf("Último preço: %.2f\n", lastPrice)
	fmt.Printf("RSI: %.2f\n", rsi)
	fmt.Println("Período:", cfg.Period)
	fmt.Println("Aberto:", IsOpened)
	fmt.Println("")

	// Obtém o estado da posição do banco
	isOpened, err := database.GetPosition(cfg.Symbol)
	if err != nil {
		log.Printf("Erro ao obter posição: %v", err)
		return
	}

	if rsi < 30 && !isOpened {
		fmt.Println("sobrevendido, momento de comprar")
		if err := NewOrder(cfg.Symbol, 0.001, "BUY", lastPrice); err != nil {
			log.Println(err)
			IsOpened = false
		} else {
			IsOpened = true
		}
	} else if rsi > 70 && isOpened {
		fmt.Println("sobrecomprado, momento de vender")
		if err := NewOrder(cfg.Symbol, 0.001, "SELL", lastPrice); err != nil {
			log.Println(err)
			IsOpened = true
		} else {
			IsOpened = false
		}
	} else {
		fmt.Println("Aguardando oportunidades...")
	}
}
