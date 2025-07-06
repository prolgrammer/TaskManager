package usecases

import (
	"TaskManager/internal/controllers/responses"
	"TaskManager/internal/repositories"
	"context"
)

type getTasksUseCase struct {
	taskRepo GetTasksRepository
}

type GetTasksUseCase interface {
	GetTasks(c context.Context, limit, offset int) ([]responses.Task, error)
}

func NewGetTasksUseCase(taskRepo repositories.TaskRepository) GetTasksUseCase {
	return &getTasksUseCase{
		taskRepo: taskRepo,
	}
}

func (g getTasksUseCase) GetTasks(c context.Context, limit, offset int) ([]responses.Task, error) {
	tasks, err := g.taskRepo.SelectAll(c, limit, offset)
	if err != nil {
		return nil, err
	}

	taskResponses := make([]responses.Task, 0, len(tasks))
	for _, task := range tasks {

		taskResponses = append(taskResponses, responses.Task{
			TaskID:    task.ID,
			Status:    string(task.Status),
			CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
			Duration:  task.Duration.String(),
		})
	}

	return taskResponses, nil
}
