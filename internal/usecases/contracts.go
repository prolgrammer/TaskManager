package usecases

import (
	"TaskManager/internal/entities"
	"context"
)

type (
	CreateTaskRepository interface {
		Insert(context context.Context, task *entities.Task) error
	}

	DeleteTaskRepository interface {
		SelectByID(ctx context.Context, id string) (entities.Task, error)
		DeleteTask(ctx context.Context, id string) error
	}

	GetTaskRepository interface {
		SelectByID(ctx context.Context, id string) (entities.Task, error)
	}

	GetTasksRepository interface {
		SelectAll(ctx context.Context, limit, offset int) ([]entities.Task, error)
	}
)
