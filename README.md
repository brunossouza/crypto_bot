# Crypto Bot

Bot de trading automatizado para criptomoedas desenvolvido em Go.

> ⚠️ **ATENÇÃO**: Este é um projeto experimental desenvolvido para fins de estudo. O código não está preparado para uso em produção e não deve ser utilizado com recursos reais.

[English version](README_EN.md)

## Funcionalidades

- Integração com a API da Binance
- Análise técnica usando RSI (Índice de Força Relativa)
- Execução automática de ordens de compra e venda
- Monitoramento em tempo real do mercado

## Requisitos

- Go 1.23 ou superior
- Conta na Binance com API Key e Secret Key

## Configuração

1. Clone o repositório
2. Copie o arquivo `.env.example` para `.env`
3. Configure as variáveis de ambiente no arquivo `.env`:

```env
API_URL=https://api.binance.com
SYMBOL=BTCUSDT
PERIOD=14
BINANCE_API_KEY=sua_api_key
BINANCE_API_SECRET=sua_api_secret
```

## Estrutura do Projeto

- `cmd/crypto_bot`: Ponto de entrada do aplicativo
- `internal/config`: Gerenciamento de configurações
- `internal/trading`: Lógica de trading e integração com a Binance

## Estratégia de Trading

O bot utiliza o indicador RSI para identificar oportunidades de mercado:
- Compra quando RSI < 30 (condição de sobrevenda)
- Vende quando RSI > 70 (condição de sobrecompra)

## Como Executar

```bash
go run cmd/crypto_bot/main.go
```

## Docker

Para executar o projeto usando Docker Compose:

```bash
# Construir e iniciar os containers
docker-compose up -d

# Verificar logs
docker-compose logs -f app

# Parar os containers
docker-compose down
```

O Docker Compose irá criar:
- Container PostgreSQL para armazenamento de dados
- Container da aplicação com todas as dependências necessárias

## Avisos

- Use por sua conta e risco
- Teste primeiro com pequenas quantias
- Monitore regularmente o funcionamento do bot
- Mantenha suas chaves API em segurança
