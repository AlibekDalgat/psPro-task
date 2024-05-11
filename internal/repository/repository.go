package repository

import (
	"github.com/jmoiron/sqlx"
	"psPro-task/internal/models"
	"time"
)

type Command interface {
	CreateCommand(command models.Command) (int, error)
	WriteResults(int, *string, *string, time.Time) error
	GetAllCommands() ([]models.Command, error)
	GetOneCommand(int) (models.Command, error)
}

type Repository struct {
	Command
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Command: NewCommandPostgres(db),
	}
}
