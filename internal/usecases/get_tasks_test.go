package usecases

import (
	"TaskManager/internal/controllers/responses"
	"TaskManager/internal/entities"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var (
	mockGetTasksTaskRepo *MockGetTasksRepository
)

func initGetTasksTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockGetTasksTaskRepo = NewMockGetTasksRepository(ctrl)
}

func TestGetTasks_Success(t *testing.T) {
	initGetTasksTestMocks(t)
	ctx := context.Background()
	limit := 10
	offset := 0

	now := time.Now()

	testTasks := []entities.Task{
		{
			ID:        "task1",
			Status:    entities.StatusCompleted,
			CreatedAt: now,
			Duration:  4 * time.Minute,
		},
		{
			ID:        "task2",
			Status:    entities.StatusRunning,
			CreatedAt: now.Add(-time.Hour),
			Duration:  0,
		},
	}

	expectedResponse := []responses.Task{
		{
			TaskID:    "task1",
			Status:    string(entities.StatusCompleted),
			CreatedAt: now.Format("2006-01-02 15:04:05"),
			Duration:  (4 * time.Minute).String(),
		},
		{
			TaskID:    "task2",
			Status:    string(entities.StatusRunning),
			CreatedAt: now.Add(-time.Hour).Format("2006-01-02 15:04:05"),
			Duration:  "0s",
		},
	}
	mockGetTasksTaskRepo.EXPECT().SelectAll(ctx, limit, offset).Return(testTasks, nil)

	useCase := NewGetTasksUseCase(mockGetTasksTaskRepo)

	response, err := useCase.GetTasks(ctx, limit, offset)

	assert.NoError(t, err)
	assert.Equal(t, expectedResponse, response)
}

func TestGetTasksUseCase_Failure_RepositoryError(t *testing.T) {
	initGetTasksTestMocks(t)
	ctx := context.Background()
	limit := 10
	offset := 0
	expectedErr := errors.New("get tasks error")

	mockGetTasksTaskRepo.EXPECT().SelectAll(ctx, limit, offset).Return(nil, expectedErr)

	useCase := NewGetTasksUseCase(mockGetTasksTaskRepo)

	result, err := useCase.GetTasks(ctx, limit, offset)

	assert.Error(t, err)
	assert.Nil(t, result)
	assert.Equal(t, expectedErr, err)
}
