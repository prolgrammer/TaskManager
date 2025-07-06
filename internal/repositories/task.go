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
	SelectByID(ctx context.Context, id string) (entities.Task, error)
	SelectAll(ctx context.Context, limit, offset int) ([]entities.Task, error)
	Count(ctx context.Context) (int, error)
	DeleteTask(ctx context.Context, id string) error
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

func (t *taskRepo) SelectByID(ctx context.Context, id string) (entities.Task, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	task, exists := t.tasks[id]
	if !exists {
		return entities.Task{}, ErrEntityNotFound
	}
	return *task, nil
}

func (t *taskRepo) SelectAll(ctx context.Context, limit, offset int) ([]entities.Task, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()

	allTasks := make([]entities.Task, 0, len(t.tasks))
	for _, task := range t.tasks {
		allTasks = append(allTasks, *task)
	}

	start := offset
	if start > len(allTasks) {
		start = len(allTasks)
	}

	end := start + limit
	if end > len(allTasks) {
		end = len(allTasks)
	}

	return allTasks[start:end], nil
}

func (t *taskRepo) Count(ctx context.Context) (int, error) {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return len(t.tasks), nil
}

func (t *taskRepo) DeleteTask(ctx context.Context, id string) error {
	t.mu.Lock()
	defer t.mu.Unlock()
	if _, exists := t.tasks[id]; !exists {
		return ErrEntityNotFound
	}
	delete(t.tasks, id)
	return nil
}
