package usecases

import (
	"TaskManager/internal/controllers/responses"
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	mockGetTaskTaskRepo *MockGetTaskRepository
)

func initGetTaskTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGetTaskTaskRepo = NewMockGetTaskRepository(ctrl)
}

func TestGetTask_Success(t *testing.T) {
	initGetTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"

	task := entities.Task{
		ID:        taskId,
		Status:    entities.StatusCompleted,
		CreatedAt: time.Now(),
		Duration:  4 * time.Minute,
	}

	expectedResponse := responses.Task{
		TaskID:    taskId,
		Status:    string(entities.StatusCompleted),
		CreatedAt: task.CreatedAt.Format("2006-01-02 15:04:05"),
		Duration:  task.Duration.String(),
	}

	mockGetTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(task, nil)

	useCase := NewGetTaskUseCase(mockGetTaskTaskRepo)

	response, err := useCase.GetTask(ctx, taskId)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGetTask_Failure_EntityNotFound(t *testing.T) {
	initGetTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"

	mockGetTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(entities.Task{}, repositories.ErrEntityNotFound)

	useCase := NewGetTaskUseCase(mockGetTaskTaskRepo)
	_, err := useCase.GetTask(ctx, taskId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)
}

func TestGetTask_Failure_RepositoryError(t *testing.T) {
	initGetTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"
	expectedErr := fmt.Errorf("get task error")

	mockGetTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(entities.Task{}, expectedErr)

	useCase := NewGetTaskUseCase(mockGetTaskTaskRepo)

	_, err := useCase.GetTask(ctx, taskId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
