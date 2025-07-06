package services

import (
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"TaskManager/pkg/logger"
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

type taskManager struct {
	logger       logger.Logger
	taskQueue    chan *entities.Task
	pendingTasks []*entities.Task
	mu           sync.Mutex
	taskRepo     repositories.TaskRepository
	stop         chan struct{}
	workersWg    sync.WaitGroup
}

type TaskManager interface {
	SubmitTask(task *entities.Task) error
	Stop()
}

func NewTaskManager(taskRepo repositories.TaskRepository, logger logger.Logger, workerCount int) TaskManager {
	tm := &taskManager{
		logger:       logger,
		taskQueue:    make(chan *entities.Task, workerCount*2),
		pendingTasks: make([]*entities.Task, 0),
		taskRepo:     taskRepo,
		stop:         make(chan struct{}),
	}

	tm.workersWg.Add(workerCount)
	for i := 0; i < workerCount; i++ {
		go tm.worker()
	}

	go tm.movePendingTasks()

	return tm
}

func (tm *taskManager) SubmitTask(task *entities.Task) error {
	select {
	case <-task.Ctx.Done():
		return task.Ctx.Err()
	default:
	}

	tm.mu.Lock()
	defer tm.mu.Unlock()

	select {
	case tm.taskQueue <- task:
		tm.logger.Debug().Msgf("Task %s submitted to taskQueue", task.ID)
		return nil
	default:
		tm.pendingTasks = append(tm.pendingTasks, task)
		tm.logger.Debug().Msgf("Task %s added to pendingTasks", task.ID)
		return nil
	}
}

func (tm *taskManager) Stop() {
	close(tm.stop)
	tm.workersWg.Wait()
}

func (tm *taskManager) worker() {
	defer tm.workersWg.Done()

	for {
		select {
		case task, ok := <-tm.taskQueue:
			if !ok {
				return
			}
			tm.processTask(task)
		case <-tm.stop:
			return
		}
	}
}

func (tm *taskManager) processTask(task *entities.Task) {
	select {
	case <-task.Ctx.Done():
		task.Status = entities.StatusFailed
		task.FinishedAt = timePtr(time.Now())
		task.Result = "Task canceled"
		_ = tm.taskRepo.Update(context.Background(), task)
		return
	default:
	}

	task.Status = entities.StatusRunning
	task.StartedAt = timePtr(time.Now())
	if err := tm.taskRepo.Update(context.Background(), task); err != nil {
		tm.logger.Error().Msgf("failed to update task: %v", err)
		return
	}

	duration := time.Duration(30+rand.Intn(20)) * time.Second
	tm.logger.Debug().Msgf("Processing task %s, duration: %v", task.ID, duration)

	select {
	case <-time.After(duration):
		task.Status = entities.StatusCompleted
		task.FinishedAt = timePtr(time.Now())
		task.Result = "Task completed successfully"
	case <-task.Ctx.Done():
		fmt.Println("ctx done")
		task.Status = entities.StatusFailed
		task.FinishedAt = timePtr(time.Now())
		task.Result = "Task canceled"
	case <-tm.stop:
		task.Status = entities.StatusFailed
		task.FinishedAt = timePtr(time.Now())
		task.Result = "Task canceled"
	}

	if err := tm.taskRepo.Update(context.Background(), task); err != nil {
		tm.logger.Error().Msgf("failed to update task: %v", err)
	}
}

func timePtr(t time.Time) *time.Time {
	return &t
}

func (tm *taskManager) movePendingTasks() {
	for {
		select {
		case <-tm.stop:
			return
		default:
			tm.mu.Lock()
			if len(tm.pendingTasks) > 0 {
				fmt.Println("trying")
				task := tm.pendingTasks[0]
				tm.pendingTasks = tm.pendingTasks[1:]
				tm.mu.Unlock()

				select {
				case tm.taskQueue <- task:
				case <-tm.stop:
					return
				}
			} else {
				tm.mu.Unlock()
				time.Sleep(time.Second)
			}
		}
	}
}
