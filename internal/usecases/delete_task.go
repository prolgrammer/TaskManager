package usecases

import (
	"TaskManager/internal/repositories"
	"TaskManager/internal/services/task_manager"
	"context"
	"errors"
)

type deleteTaskUseCase struct {
	taskRepo    repositories.TaskRepository
	taskManager task_manager.TaskManager
}

type DeleteTaskUseCase interface {
	DeleteTask(c context.Context, taskId string) error
}

func NewDeleteTaskUseCase(
	taskRepo repositories.TaskRepository,
	taskManager task_manager.TaskManager,
) DeleteTaskUseCase {
	return &deleteTaskUseCase{
		taskRepo:    taskRepo,
		taskManager: taskManager,
	}
}

func (d *deleteTaskUseCase) DeleteTask(c context.Context, taskId string) error {
	task, err := d.taskRepo.SelectByID(c, taskId)
	if err != nil {
		if errors.Is(err, repositories.ErrEntityNotFound) {
			return ErrEntityNotFound
		}
		return err
	}

	if task.FinishedAt != nil {
		return d.taskRepo.Delete(c, taskId)
	}
	d.taskManager.CancelTask(taskId)

	return d.taskRepo.Delete(c, taskId)
}
