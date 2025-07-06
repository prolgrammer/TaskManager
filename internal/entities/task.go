package entities

import (
	"context"
	"time"
)

type TaskStatus string

const (
	StatusCreated   TaskStatus = "created"
	StatusRunning   TaskStatus = "running"
	StatusCompleted TaskStatus = "completed"
	StatusFailed    TaskStatus = "failed"
)

type Task struct {
	ID         string
	Text       string
	Status     TaskStatus
	CreatedAt  time.Time
	StartedAt  *time.Time
	FinishedAt *time.Time
	Duration   time.Duration
	Result     string
	Error      string
	Ctx        context.Context
}
