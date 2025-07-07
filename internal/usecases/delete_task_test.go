package usecases

import (
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mockDeleteTaskTaskRepo    *MockDeleteTaskRepository
	mockDeleteTaskTaskManager *MockDeleteTaskTaskManager
)

func initDeleteTaskTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockDeleteTaskTaskRepo = NewMockDeleteTaskRepository(ctrl)
	mockDeleteTaskTaskManager = NewMockDeleteTaskTaskManager(ctrl)
}

func TestDeleteTask_Success(t *testing.T) {
	initDeleteTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"

	task := entities.Task{ID: taskId, FinishedAt: nil}

	mockDeleteTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(task, nil)
	mockDeleteTaskTaskManager.EXPECT().CancelTask(taskId)
	mockDeleteTaskTaskRepo.EXPECT().DeleteTask(ctx, taskId).Return(nil)

	useCase := NewDeleteTaskUseCase(mockDeleteTaskTaskRepo, mockDeleteTaskTaskManager)

	err := useCase.DeleteTask(ctx, taskId)

	assert.NoError(t, err)
}

func TestDeleteTask_Failure_EntityNotFound(t *testing.T) {
	initDeleteTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"

	mockDeleteTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(entities.Task{}, repositories.ErrEntityNotFound)

	useCase := NewDeleteTaskUseCase(mockDeleteTaskTaskRepo, mockDeleteTaskTaskManager)

	err := useCase.DeleteTask(ctx, taskId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityNotFound)

}

func TestDeleteTask_Failure_DeleteFinishedTask(t *testing.T) {
	initDeleteTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"
	expectedError := errors.New("select error")
	mockDeleteTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(entities.Task{}, expectedError)

	useCase := NewDeleteTaskUseCase(mockDeleteTaskTaskRepo, mockDeleteTaskTaskManager)
	err := useCase.DeleteTask(ctx, taskId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}

func TestDeleteTask_Failure_DeleteUnfinishedTask(t *testing.T) {
	initDeleteTaskTestMocks(t)
	ctx := context.Background()
	taskId := "test"
	expectedError := errors.New("delete failed")
	task := entities.Task{ID: taskId, FinishedAt: nil}

	mockDeleteTaskTaskRepo.EXPECT().SelectByID(ctx, taskId).Return(task, nil)
	mockDeleteTaskTaskManager.EXPECT().CancelTask(taskId)
	mockDeleteTaskTaskRepo.EXPECT().DeleteTask(ctx, taskId).Return(expectedError)

	useCase := NewDeleteTaskUseCase(mockDeleteTaskTaskRepo, mockDeleteTaskTaskManager)

	err := useCase.DeleteTask(ctx, taskId)
	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedError)
}
