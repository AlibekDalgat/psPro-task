package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"psPro-task/internal/config"
)

const (
	commandsTable = "commands"
)

func OpenDB(cfg config.DBConfig) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.User, cfg.Database, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}
