package usecases

import (
	"TaskManager/internal/entities"
	"context"
)

type (
	CreateTaskRepository interface {
		Insert(context context.Context, task *entities.Task) error
	}
)
