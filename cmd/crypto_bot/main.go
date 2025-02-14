package main

import (
	"fmt"
	"log"
	"time"

	"github.com/brunossouza/crypto_bot/internal/config"
	"github.com/brunossouza/crypto_bot/internal/trading"
)

// main é o ponto de entrada principal do bot de criptomoedas.
// Responsabilidades:
// - Carregar as configurações do sistema através do arquivo .env
// - Inicializar o módulo de trading com as configurações carregadas
// - Configurar um temporizador para execução periódica das operações
// - Manter o bot em execução contínua até ser manualmente interrompido
//
// Fluxo de execução:
// 1. Carrega as configurações do arquivo .env
// 2. Inicializa o módulo de trading
// 3. Configura um temporizador de 10 segundos
// 4. Executa a primeira operação de trading
// 5. Entra em loop infinito, executando operações a cada 10 segundos
//
// Em caso de erro na carga das configurações, o programa é encerrado com log.Fatal
func main() {
	conf, err := config.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	// Inicializa o pacote de trading com as configurações
	trading.Initialize(conf)

	// Cria um temporizador que dispara a cada 10 segundos
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	fmt.Println("Bot iniciado! Pressione CTRL+C para parar")
	trading.StartTrading()

	// Executa indefinidamente até o programa ser interrompido
	for range ticker.C {
		trading.StartTrading()
	}
}
