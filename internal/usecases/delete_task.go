package usecases

import (
	"TaskManager/internal/repositories"
	"context"
	"errors"
)

type deleteTaskUseCase struct {
	taskRepo repositories.TaskRepository
}

type DeleteTaskUseCase interface {
	DeleteTask(c context.Context, taskId string) error
}

func NewDeleteTaskUseCase(
	taskRepo repositories.TaskRepository,
) DeleteTaskUseCase {
	return &deleteTaskUseCase{
		taskRepo: taskRepo,
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

	task.Ctx.Done()
	//if err = d.taskRepo.Delete(c, taskId); err != nil {
	//	return err
	//}

	return nil
}
