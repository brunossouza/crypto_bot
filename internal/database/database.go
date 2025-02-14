package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

type Order struct {
	ID        int64     `json:"id"`
	Symbol    string    `json:"symbol"`
	Side      string    `json:"side"`
	Quantity  float64   `json:"quantity"`
	Price     float64   `json:"price"`
	CreatedAt time.Time `json:"created_at"`
}

type Position struct {
	ID        int64     `json:"id"`
	Symbol    string    `json:"symbol"`
	IsOpened  bool      `json:"is_opened"`
	UpdatedAt time.Time `json:"updated_at"`
}

var db *sql.DB

func Initialize() error {
	// Cria a connection string a partir das variáveis de ambiente
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

	// Criar tabelas se não existirem
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
	_, err = db.Exec(createTables)
	return err
}

func SaveOrder(symbol string, side string, quantity, price float64) error {
	stmt, err := db.Prepare(`
		INSERT INTO orders (symbol, side, quantity, price)
		VALUES ($1, $2, $3, $4)
	`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	_, err = stmt.Exec(symbol, side, quantity, price)
	return err
}

func UpdatePosition(symbol string, isOpened bool) error {
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

	_, err = stmt.Exec(symbol, isOpened, isOpened)
	return err
}

func GetPosition(symbol string) (bool, error) {
	var isOpened bool
	err := db.QueryRow(`
		SELECT is_opened FROM positions
		WHERE symbol = $1
	`, symbol).Scan(&isOpened)

	if err == sql.ErrNoRows {
		return false, nil
	}
	return isOpened, err
}

func Close() {
	if db != nil {
		db.Close()
	}
}
