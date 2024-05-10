package config

import "os"

type DBConfig struct {
	User     string
	Password string
	Host     string
	Port     string
	Database string
	SSLMode  string
}

type ServerConfig struct {
	HTTPPort string
}

func GetDBConfig() DBConfig {
	return DBConfig{
		User:     os.Getenv("PGUSER"),
		Password: os.Getenv("PGPASSWORD"),
		Host:     os.Getenv("PGHOST"),
		Port:     os.Getenv("PGPORT"),
		Database: os.Getenv("PGDATABASE"),
		SSLMode:  os.Getenv("PGSSLMODE"),
	}
}
