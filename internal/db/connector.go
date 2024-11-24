package db

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/VadimBorzenkov/gw-currency-wallet/internal/config"
)

// Init инициализирует подключение к базе данных на основе конфигурации и возвращает объект базы данных.
func Init(cfg *config.Config) (*sql.DB, error) {
	url := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPass, cfg.DBName)

	db, err := sql.Open("postgres", url)
	if err != nil {
		log.Printf("Could not connect to DB: %v", err)
		return nil, err
	}
	return db, nil
}

// Close закрывает соединение с базой данных и возвращает ошибку, если возникла проблема при закрытии.
func Close(db *sql.DB) error {
	if err := db.Close(); err != nil {
		return err
	}
	return nil
}
