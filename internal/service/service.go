package service

import (
	"psPro-task/internal/models"
	"psPro-task/internal/repository"
)

type Command interface {
	CreateCommand(command models.Command) (int, error)
	ExecuteCommand(int, string)
	StopCommand(int) error
	StartCommand(int) error
	KillCommand(int) error
}

type Service struct {
	Command
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Command: NewCommandService(repo),
	}
}
