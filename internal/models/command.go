package models

import (
	"time"
)

type Command struct {
	Id         int        `json:"id" db:"id"`
	Script     string     `json:"script" db:"script"`
	CreatedAt  time.Time  `json:"created_at" db:"created_at"`
	ExecutedAT *time.Time `json:"executed_at" db:"executed_at"`
	Stdout     *string    `json:"stdout" db:"stdout"`
	Stderr     *string    `json:"stderr" db:"stderr"`
}

type Job struct {
	ActionChan chan string
	IsRun      bool
}
