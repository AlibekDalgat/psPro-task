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
	query := fmt.Sprintf("INSERT INTO %s (script, created_at) VALUES ($1, $2) RETURNING id", commandsTable)
	row := p.db.QueryRow(query, command.Script, command.CreatedAt)
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (p *CommandPostgres) AddRecordToColumn(column string, id int, data string) error {
	query := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", commandsTable, column)
	_, err := p.db.Exec(query, data, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *CommandPostgres) RewriteColumn(column string, id int, data interface{}) error {
	query := fmt.Sprintf("UPDATE %s SET %s = $1 WHERE id = $2", commandsTable, column)
	_, err := p.db.Exec(query, data, id)
	if err != nil {
		return err
	}
	return nil
}

func (p *CommandPostgres) GetAllCommands() ([]models.Command, error) {
	var commands []models.Command
	query := fmt.Sprintf("SELECT * FROM %s", commandsTable)
	err := p.db.Select(&commands, query)
	if err != nil {
		return nil, err
	}
	return commands, nil
}

func (p *CommandPostgres) GetOneCommand(id int) (models.Command, error) {
	var command models.Command
	query := fmt.Sprintf("SELECT * FROM %s WHERE id = $1", commandsTable)
	row := p.db.QueryRow(query, id)
	if err := row.Scan(&command); err != nil {
		return models.Command{}, err
	}
	return command, nil
}
