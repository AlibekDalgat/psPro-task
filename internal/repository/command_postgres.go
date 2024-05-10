package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"psPro-task/internal/models"
)

type CommandPostgres struct {
	db *sqlx.DB
}

func NewCommandPostgres(db *sqlx.DB) *CommandPostgres {
	return &CommandPostgres{db: db}
}

func (p *CommandPostgres) CreateCommand(command models.Command) (int, error) {
	var id int
	query := fmt.Sprintf("INSERT INTO %s (script, created_at) VALUES ($1, $2)", commandsTable)
	row := p.db.QueryRow(query, command.Script, command.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *CommandPostgres) WriteToColumn(column string, id int, dataa interface{}) error {
	//TODO implement me
	panic("implement me")
}
