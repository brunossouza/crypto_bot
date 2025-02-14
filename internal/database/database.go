package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// Order representa um registro de ordem no sistema.
// Cada ordem contém informações sobre o símbolo, o tipo de operação (Buy/Sell),
// quantidade, preço e a data de criação.
type Order struct {
	// ID é o identificador único da ordem.
	ID int64 `json:"id"`
	// Symbol é o ativo negociado.
	Symbol string `json:"symbol"`
	// Side indica se a ordem é de compra ou venda.
	Side string `json:"side"`
	// Quantity representa a quantidade negociada.
	Quantity float64 `json:"quantity"`
	// Price é o valor da ordem.
	Price float64 `json:"price"`
	// CreatedAt marca o momento em que a ordem foi criada.
	CreatedAt time.Time `json:"created_at"`
}

// Position representa a posição atual para um símbolo.
// Armazena se a posição está aberta e a data da última atualização.
type Position struct {
	// ID é o identificador único da posição.
	ID int64 `json:"id"`
	// Symbol é o ativo relacionado à posição.
	Symbol string `json:"symbol"`
	// IsOpened indica se a posição está atualmente aberta (true) ou fechada (false).
	IsOpened bool `json:"is_opened"`
	// UpdatedAt indica o momento da última atualização da posição.
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB

// Initialize estabelece a conexão com o banco de dados PostgreSQL e inicializa as tabelas necessárias.
// Utiliza as seguintes variáveis de ambiente para configuração:
//   - DB_HOST: endereço do servidor PostgreSQL
//   - DB_PORT: porta do servidor PostgreSQL
//   - DB_USER: usuário do banco de dados
//   - DB_PASSWORD: senha do usuário
//   - DB_NAME: nome do banco de dados
//
// Retorna erro se:
//   - Falhar ao estabelecer conexão com o banco
//   - Falhar ao criar as tabelas necessárias
func Initialize() error {
	// Cria a connection string a partir das variáveis de ambiente.
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	var err error
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		return err
	}

	// Declaração das queries de criação das tabelas.
	createTables := `
	CREATE TABLE IF NOT EXISTS orders (
		id SERIAL PRIMARY KEY,
		symbol TEXT NOT NULL,
		side TEXT NOT NULL,
		quantity REAL NOT NULL,
		price REAL NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS positions (
		id SERIAL PRIMARY KEY,
		symbol TEXT UNIQUE NOT NULL,
		is_opened BOOLEAN NOT NULL,
		updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	`
	// Executa as queries para criar as tabelas se estas ainda não existirem.
	_, err = db.Exec(createTables)
	return err
}

// SaveOrder registra uma nova ordem de compra ou venda no banco de dados.
// Parâmetros:
//   - symbol: identificador do par de moedas (ex: "BTCUSDT")
//   - side: direção da ordem ("BUY" ou "SELL")
//   - quantity: quantidade do ativo a ser negociado
//   - price: preço unitário do ativo
//
// Retorna erro se:
//   - Falhar ao preparar a declaração SQL
//   - Falhar ao executar a inserção no banco
func SaveOrder(symbol string, side string, quantity, price float64) error {
	// Prepara a instrução SQL para inserir a ordem.
	stmt, err := db.Prepare(`
		INSERT INTO orders (symbol, side, quantity, price)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Executa a instrução com os parâmetros passados.
	_, err = stmt.Exec(symbol, side, quantity, price)
	return err
}

// UpdatePosition atualiza ou cria uma nova posição para um determinado símbolo no banco de dados.
// Utiliza a cláusula ON CONFLICT para garantir que existe apenas uma posição por símbolo.
// Parâmetros:
//   - symbol: identificador do par de moedas (ex: "BTCUSDT")
//   - isOpened: true se a posição está aberta, false se fechada
//
// Retorna erro se:
//   - Falhar ao preparar a declaração SQL
//   - Falhar ao executar a atualização/inserção no banco
func UpdatePosition(symbol string, isOpened bool) error {
	// Prepara a instrução SQL que insere uma nova posição ou atualiza a existente.
	stmt, err := db.Prepare(`
		INSERT INTO positions (symbol, is_opened, updated_at)
		VALUES ($1, $2, CURRENT_TIMESTAMP)
		ON CONFLICT (symbol)
		DO UPDATE SET is_opened = $3, updated_at = CURRENT_TIMESTAMP
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	// Executa a instrução com os parâmetros fornecidos.
	_, err = stmt.Exec(symbol, isOpened, isOpened)
	return err
}

// GetPosition consulta o estado atual da posição para um determinado símbolo.
// Parâmetros:
//   - symbol: identificador do par de moedas (ex: "BTCUSDT")
//
// Retorna:
//   - bool: true se existe uma posição aberta, false se fechada ou inexistente
//   - error: erro em caso de falha na consulta ao banco
//
// Comportamento especial:
//   - Se não existir posição para o símbolo, retorna (false, nil)
//   - Se ocorrer erro na consulta, retorna (false, erro)
func GetPosition(symbol string) (bool, error) {
	var isOpened bool
	// Executa a consulta e mapeia o resultado para a variável isOpened.
	err := db.QueryRow(`
		SELECT is_opened FROM positions
		WHERE symbol = $1
	`, symbol).Scan(&isOpened)

	// Se não houver linha, retorna false sem erro.
	if err == sql.ErrNoRows {
		return false, nil
	}
	return isOpened, err
}

// Close finaliza a conexão com o banco de dados de forma segura.
// Deve ser chamado quando a aplicação for encerrada para liberar recursos.
// É seguro chamar mesmo se a conexão não estiver inicializada (db == nil).
func Close() {
	if db != nil {
		db.Close()
	}
}
