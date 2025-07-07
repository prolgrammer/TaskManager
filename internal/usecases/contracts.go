package usecases

import (
	"TaskManager/internal/entities"
	"context"
)

//go:generate mockgen -source=contracts.go --destination=mock_test.go -package=usecases

type (
	CreateTaskTaskManager interface {
		SubmitTask(task *entities.Task) error
	}

	CreateTaskRepository interface {
		Insert(context context.Context, task *entities.Task) error
	}

	DeleteTaskRepository interface {
		SelectByID(ctx context.Context, id string) (entities.Task, error)
		DeleteTask(ctx context.Context, id string) error
	}

	DeleteTaskTaskManager interface {
		CancelTask(taskID string)
	}

	GetTaskRepository interface {
		SelectByID(ctx context.Context, id string) (entities.Task, error)
	}

	GetTasksRepository interface {
		SelectAll(ctx context.Context, limit, offset int) ([]entities.Task, error)
	}
)
