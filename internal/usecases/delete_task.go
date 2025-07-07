package usecases

import (
	"TaskManager/internal/repositories"
	"context"
	"errors"
)

type deleteTaskUseCase struct {
	taskRepo    DeleteTaskRepository
	taskManager DeleteTaskTaskManager
}

type DeleteTaskUseCase interface {
	DeleteTask(c context.Context, taskId string) error
}

func NewDeleteTaskUseCase(
	taskRepo DeleteTaskRepository,
	taskManager DeleteTaskTaskManager,
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
		return d.taskRepo.DeleteTask(c, taskId)
	}
	d.taskManager.CancelTask(taskId)

	return d.taskRepo.DeleteTask(c, taskId)
}
