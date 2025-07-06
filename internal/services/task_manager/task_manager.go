package task_manager

import (
	"TaskManager/internal/entities"
	"TaskManager/internal/repositories"
	"TaskManager/pkg/logger"
	"context"
	"sync"
	"time"
)

type TaskManager interface {
	SubmitTask(task *entities.Task) error
	CancelTask(taskID string)
	Stop()
}

type taskManager struct {
	logger       logger.Logger
	mu           sync.Mutex
	workersWg    sync.WaitGroup
	taskQueue    chan *entities.Task
	stop         chan struct{}
	activeTasks  map[string]context.CancelFunc
	pendingTasks []*entities.Task
	taskRepo     repositories.TaskRepository
	workerPool   *workerPool
}

func NewTaskManager(taskRepo repositories.TaskRepository, logger logger.Logger, workerCount int) TaskManager {
	tm := &taskManager{
		logger:       logger,
		taskQueue:    make(chan *entities.Task, workerCount*2),
		stop:         make(chan struct{}),
		activeTasks:  make(map[string]context.CancelFunc),
		pendingTasks: make([]*entities.Task, 0),
		taskRepo:     taskRepo,
	}

	tm.workerPool = newWorkerPool(workerCount, tm)
	go tm.movePendingTasks()

	return tm
}

func (tm *taskManager) SubmitTask(task *entities.Task) error {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	ctx, cancel := context.WithCancel(context.Background())
	task.Ctx = ctx
	tm.activeTasks[task.ID] = cancel

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

func (tm *taskManager) CancelTask(taskID string) {
	tm.mu.Lock()
	defer tm.mu.Unlock()

	cancel, exists := tm.activeTasks[taskID]
	if exists {
		cancel()
		delete(tm.activeTasks, taskID)
		return
	}

	for i, task := range tm.pendingTasks {
		if task.ID == taskID {
			tm.pendingTasks = append(tm.pendingTasks[:i], tm.pendingTasks[i+1:]...)
			return
		}
	}

	return
}

func (tm *taskManager) Stop() {
	close(tm.stop)
	tm.workerPool.wait()
}

func (tm *taskManager) movePendingTasks() {
	for {
		select {
		case <-tm.stop:
			return
		default:
			tm.mu.Lock()
			if len(tm.pendingTasks) > 0 {
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
