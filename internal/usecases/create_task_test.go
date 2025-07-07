package usecases

import (
	"TaskManager/internal/controllers/requests"
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"testing"
)

var (
	mockCreateTaskTaskRepo    *MockCreateTaskRepository
	mockCreateTaskTaskManager *MockCreateTaskTaskManager
)

func initCreateTaskTestMocks(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCreateTaskTaskRepo = NewMockCreateTaskRepository(ctrl)
	mockCreateTaskTaskManager = NewMockCreateTaskTaskManager(ctrl)
}

func TestCreateTask_Success(t *testing.T) {
	initCreateTaskTestMocks(t)
	ctx := context.Background()
	req := requests.CreateTask{
		Text: "test",
	}

	task := gomock.AssignableToTypeOf(&entities.Task{})

	mockCreateTaskTaskRepo.EXPECT().Insert(ctx, task).Return(nil)
	mockCreateTaskTaskManager.EXPECT().SubmitTask(task).Return(nil)

	useCase := NewCreateTaskUseCase(mockCreateTaskTaskRepo, mockCreateTaskTaskManager)

	response, err := useCase.CreateTask(ctx, req)

	assert.NoError(t, err)
	assert.NotNil(t, response)
}

func TestCreateTask_Failure_EntityAlreadyExists(t *testing.T) {
	initCreateTaskTestMocks(t)
	ctx := context.Background()
	req := requests.CreateTask{
		Text: "test",
	}

	mockCreateTaskTaskRepo.EXPECT().Insert(ctx, gomock.Any()).Return(repositories.ErrEntityAlreadyExists)

	useCase := NewCreateTaskUseCase(mockCreateTaskTaskRepo, mockCreateTaskTaskManager)

	_, err := useCase.CreateTask(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, ErrEntityAlreadyExists)
}

func TestCreateTask_Failure_Insert(t *testing.T) {
	initCreateTaskTestMocks(t)
	ctx := context.Background()
	req := requests.CreateTask{
		Text: "test",
	}

	expectedErr := fmt.Errorf("insert error")
	mockCreateTaskTaskRepo.EXPECT().Insert(ctx, gomock.Any()).Return(expectedErr)

	useCase := NewCreateTaskUseCase(mockCreateTaskTaskRepo, mockCreateTaskTaskManager)
	_, err := useCase.CreateTask(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}

func TestCreateTask_Failure_SubmitError(t *testing.T) {
	initCreateTaskTestMocks(t)
	ctx := context.Background()
	req := requests.CreateTask{
		Text: "test",
	}

	expectedErr := fmt.Errorf("task manager error")

	mockCreateTaskTaskRepo.EXPECT().Insert(ctx, gomock.Any()).Return(nil)
	mockCreateTaskTaskManager.EXPECT().SubmitTask(gomock.Any()).Return(expectedErr)

	useCase := NewCreateTaskUseCase(mockCreateTaskTaskRepo, mockCreateTaskTaskManager)

	_, err := useCase.CreateTask(ctx, req)

	assert.Error(t, err)
	assert.ErrorIs(t, err, expectedErr)
}
