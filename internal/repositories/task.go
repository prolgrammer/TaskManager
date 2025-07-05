package repositories

import (
	"TaskManager/internal/entities"
	"context"
	"sync"
)

type taskRepo struct {
	tasks map[string]*entities.Task
	mu    sync.RWMutex
}

type TaskRepository interface {
	Insert(context context.Context, task *entities.Task) error
	Update(context context.Context, task *entities.Task) error
}

func NewTaskRepository() TaskRepository {
	return &taskRepo{
		tasks: make(map[string]*entities.Task),
	}
}

func (t *taskRepo) Insert(context context.Context, task *entities.Task) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.tasks[task.ID]; exists {
		return ErrEntityAlreadyExists
	}

	t.tasks[task.ID] = task
	return nil
}

func (t *taskRepo) Update(context context.Context, task *entities.Task) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.tasks[task.ID]; !exists {
		return ErrEntityNotFound
	}
	t.tasks[task.ID] = task
	return nil
}
