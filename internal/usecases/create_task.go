package usecases

import (
	"TaskManager/internal/controllers/requests"
	"TaskManager/internal/controllers/responses"
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"TaskManager/internal/services"
	"context"
	"errors"
	"fmt"
	"github.com/google/uuid"
	"time"
)

type createTaskUseCase struct {
	taskRepo    CreateTaskRepository
	taskManager services.TaskManager
}

type CreateTaskUseCase interface {
	CreateTask(context context.Context, req requests.CreateTask) (responses.Task, error)
}

func NewCreateTaskUseCase(taskRepo CreateTaskRepository, taskManager services.TaskManager) CreateTaskUseCase {
	return &createTaskUseCase{
		taskRepo:    taskRepo,
		taskManager: taskManager,
	}
}

func (c *createTaskUseCase) CreateTask(context context.Context, req requests.CreateTask) (responses.Task, error) {
	task := &entities.Task{
		ID:        uuid.New().String(),
		Text:      req.Text,
		Status:    entities.StatusCreated,
		CreatedAt: time.Now(),
		Ctx:       context,
	}

	if err := c.taskRepo.Insert(context, task); err != nil {
		if errors.Is(err, repositories.ErrEntityAlreadyExists) {
			return responses.Task{}, ErrEntityAlreadyExists
		}
		return responses.Task{}, fmt.Errorf("failed to create task: %w", err)
	}

	if err := c.taskManager.SubmitTask(task); err != nil {
		return responses.Task{}, fmt.Errorf("failed to submit task: %w", err)
	}

	return responses.Task{
		TaskID:    task.ID,
		Status:    string(task.Status),
		CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
	}, nil
}
