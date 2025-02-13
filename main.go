package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"strconv"
	"time"
)

const (
	// API_URL is the base URL for Binance API
	API_URL = "https://testnet.binance.vision" //"https://api.binance.com"
	// SYMBOL is the currency pair to trade
	SYMBOL = "BTCUSDT"
	// PERIOD is number of candles to calculate the average - used for RSI calculation
	PERIOD = 14
)

var (
	// IsOpened is a flag to check if the position is opened
	IsOpened bool = false
)

type Candlestick struct {
	OpenTime                 int64   // Kline open time in milliseconds
	Open                     float64 // Opening price of the candle
	High                     float64 // Highest price during the candle period
	Low                      float64 // Lowest price during the candle period
	Close                    float64 // Closing price (or latest price) of the candle
	Volume                   float64 // Total trading volume during the candle period
	CloseTime                int64   // Kline close time in milliseconds
	QuoteAssetVolume         float64 // Total quote asset volume during the candle period
	NumberOfTrades           int64   // Number of trades during the candle period
	TakerBuyBaseAssetVolume  float64 // Taker buy base asset volume
	TakerBuyQuoteAssetVolume float64 // Taker buy quote asset volume
	Ignore                   float64 // Unused field, ignore
}

func GetCandlesticks(symbol string, interval string, limit int) []Candlestick {
	// Create a new request
	req, err := http.NewRequest("GET", fmt.Sprintf("%s/api/v3/klines?symbol=%s&interval=%s&limit=%d", API_URL, symbol, interval, limit), nil)
	if err != nil {
		log.Fatal(err)
	}

	// Send the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}

	// Parse the JSON response
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

// Helper function to parse string to float64
func parseFloat(str string) float64 {
	val, err := strconv.ParseFloat(str, 64)
	if err != nil {
		log.Fatal(err)
	}
	return val
}

// calculateAverage calculates the average gains and losses for a given period
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

// CalculateRSI calculates the Relative Strength Index using EMA
func CalculateRSI(prices []float64, period int) float64 {
	if len(prices) < period+1 {
		log.Fatal("Not enough prices to calculate RSI")
	}

	var avgGains, avgLoss float64

	// Calculate initial values
	for i := 1; i < len(prices); i++ {
		gain, loss := calculateAverage(prices, period, i)

		if i == 1 {
			avgGains = gain
			avgLoss = loss
			continue
		}

		// Calculate EMA(Exponential Moving Average) for gains and losses
		avgGains = (avgGains*float64(period-1) + gain) / float64(period)
		avgLoss = (avgLoss*float64(period-1) + loss) / float64(period)
	}

	rs := avgGains / avgLoss
	rsi := 100.0 - (100.0 / (1.0 + rs))

	return rsi
}

func StartTrading() {
	// Get candlesticks data
	candlesticks := GetCandlesticks(SYMBOL, "15m", 100)

	// Get the last price
	lastPrice := candlesticks[len(candlesticks)-1].Close

	var prices []float64
	for _, candlestick := range candlesticks {
		prices = append(prices, candlestick.Close)
	}

	// Calculate the RSI
	rsi := CalculateRSI(prices, PERIOD)

	// Clear the screen
	fmt.Print("\033[H\033[2J")

	// Print the last price
	fmt.Printf("Last price: %.2f\n", lastPrice)
	// Print the RSI
	fmt.Printf("RSI: %.2f\n", rsi)
}

func main() {
	// Create a ticker that ticks every 3 seconds
	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	fmt.Println("Bot started! Press CTRL+C to stop")

	// Run forever until program is interrupted
	for range ticker.C {
		StartTrading()
	}
}
