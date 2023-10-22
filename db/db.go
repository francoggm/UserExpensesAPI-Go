package db

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/francoggm/go_expenses_api/configs"

	_ "github.com/lib/pq"
)

func createTables(db *sql.DB) error {
	_, err := db.Exec(`
	CREATE TABLE IF NOT EXISTS users (
		id SERIAL PRIMARY KEY,
		email VARCHAR(255) NOT NULL,
		name VARCHAR(255) NOT NULL,
		password BYTEA,
		created_at TIMESTAMP NOT NULL,
		last_login TIMESTAMP
	);

	CREATE TABLE IF NOT EXISTS expenses (
		id SERIAL PRIMARY KEY,
		user_id INT NOT NULL,
		title VARCHAR(150) NOT NULL,
		description TEXT,
		value NUMERIC(12, 2) NOT NULL,
		category_type SMALLINT,
		movimentation_type SMALLINT,
		created_at TIMESTAMP NOT NULL,
		CONSTRAINT fk_users FOREIGN KEY(user_id) REFERENCES users(id)
	);
	`)

	return err
}

func NewDatabase() (*sql.DB, error) {
	cfg := configs.GetConfigs()

	// deploy env config
	sc := os.Getenv("DATABASE_URL")
	if sc == "" {
		sc = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DB)
	}

	db, err := sql.Open("postgres", sc)
	if err != nil {
		return nil, err
	}

	err = createTables(db)
	if err != nil {
		fmt.Print("error creating tables!")
	}

	return db, nil
}
