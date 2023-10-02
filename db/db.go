package db

import (
	"database/sql"
	"expenses_api/configs"
	"fmt"

	_ "github.com/lib/pq"
)

func NewDatabase() (*sql.DB, error) {
	cfg := configs.GetConfigs()

	sc := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DB)

	db, err := sql.Open("postgres", sc)
	if err != nil {
		return nil, err
	}

	return db, nil
}
