package usecases

import (
	"TaskManager/internal/controllers/responses"
	"TaskManager/internal/repositories"
	"context"
	"errors"
)

type getTaskUseCase struct {
	taskRepo GetTaskRepository
}

type GetTaskUseCase interface {
	GetTask(context context.Context, id string) (responses.Task, error)
}

func NewGetTaskUseCase(taskRepo repositories.TaskRepository) GetTaskUseCase {
	return &getTaskUseCase{
		taskRepo: taskRepo,
	}
}

func (g *getTaskUseCase) GetTask(context context.Context, id string) (responses.Task, error) {
	task, err := g.taskRepo.SelectByID(context, id)
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return responses.Task{}, ErrEntityNotFound
		}
		return responses.Task{}, err
	}

	return responses.Task{
		TaskID:    task.ID,
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		Duration:  task.Duration.String(),
	}, nil
}
