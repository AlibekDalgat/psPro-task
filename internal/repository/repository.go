package repository

import (
	"github.com/jmoiron/sqlx"
	"psPro-task/internal/models"
)

type Command interface {
	CreateCommand(command models.Command) (int, error)
	WriteToColumn(string, int, interface{}) error
}

type Repository struct {
	Command
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Command: NewCommandPostgres(db),
	}
}
