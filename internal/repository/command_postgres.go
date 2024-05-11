package repository

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"psPro-task/internal/models"
	"time"
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

func (p *CommandPostgres) WriteResults(id int, stdout, stderr *string, executedAt time.Time) error {
	if *stdout == "" {
		stdout = nil
	}
	if *stderr == "" {
		stderr = nil
	}
	query := fmt.Sprintf("UPDATE %s SET stdout = $1, stderr = $2, executed_at = $3 WHERE id = $4", commandsTable)
	_, err := p.db.Exec(query, stdout, stderr, executedAt, id)
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
	query := fmt.Sprintf("SELECT id, script, created_at, executed_at, stdout, stderr FROM %s WHERE id = $1",
		commandsTable)
	row := p.db.QueryRow(query, id)
	if err := row.Scan(&command.Id, &command.Script, &command.CreatedAt, &command.ExecutedAT, &command.Stdout,
		&command.Stderr); err != nil {
		return models.Command{}, err
	}
	return command, nil
}
