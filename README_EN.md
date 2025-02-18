# Crypto Bot

Automated cryptocurrency trading bot developed in Go.

> ⚠️ **WARNING**: This is an experimental project developed for study purposes. The code is not production-ready and should not be used with real assets.

## Features

- Binance API integration
- Technical analysis using RSI (Relative Strength Index)
- Automatic buy and sell order execution
- Real-time market monitoring

## Requirements

- Go 1.23 or higher
- Binance account with API Key and Secret Key

## Setup

1. Clone the repository
2. Copy `.env.example` to `.env`
3. Configure environment variables in the `.env` file:

```env
API_URL=https://api.binance.com
SYMBOL=BTCUSDT
PERIOD=14
BINANCE_API_KEY=your_api_key
BINANCE_API_SECRET=your_api_secret
```

## Project Structure

- `cmd/crypto_bot`: Application entry point
- `internal/config`: Configuration management
- `internal/trading`: Trading logic and Binance integration

## Trading Strategy

The bot uses the RSI indicator to identify market opportunities:
- Buys when RSI < 30 (oversold condition)
- Sells when RSI > 70 (overbought condition)

## How to Run

```bash
go run cmd/crypto_bot/main.go
```

## Docker

To run the project using Docker Compose:

```bash
# Build and start containers
docker-compose up -d

# Check logs
docker-compose logs -f app

# Stop containers
docker-compose down
```

Docker Compose will create:
- PostgreSQL container for data storage
- Application container with all necessary dependencies

## Disclaimers

- Use at your own risk
- Test first with small amounts
- Regularly monitor the bot's operation
- Keep your API keys secure
