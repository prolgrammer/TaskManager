package task_manager

import (
	"TaskManager/internal/entities"
	"context"
	"math/rand"
	"sync"
	"time"
)

type workerPool struct {
	wg sync.WaitGroup
	tm *taskManager
}

func newWorkerPool(count int, tm *taskManager) *workerPool {
	wp := &workerPool{
		tm: tm,
	}
	wp.wg.Add(count)
	for i := 0; i < count; i++ {
		go wp.worker()
	}
	return wp
}

func (wp *workerPool) worker() {
	defer wp.wg.Done()

	for {
		select {
		case task, ok := <-wp.tm.taskQueue:
			if !ok {
				return
			}
			wp.processTask(task)
		case <-wp.tm.stop:
			return
		}
	}
}

func (wp *workerPool) processTask(task *entities.Task) {
	startTime := time.Now()
	defer func() {
		wp.tm.mu.Lock()
		delete(wp.tm.activeTasks, task.ID)
		wp.tm.mu.Unlock()
	}()

	select {
	case <-task.Ctx.Done():
		wp.tm.logger.Debug().Msgf("Task %s canceled before processing", task.ID)
		return
	default:
	}

	task.Status = entities.StatusRunning
	task.StartedAt = timePtr(startTime)
	if err := wp.tm.taskRepo.Update(context.Background(), task); err != nil {
		wp.tm.logger.Error().Msgf("failed to update task: %v", err)
		return
	}

	duration := time.Duration(3+rand.Intn(2)) * time.Minute
	wp.tm.logger.Debug().Msgf("Processing task %s, duration: %v", task.ID, duration)

	select {
	case <-time.After(duration):
		wp.finalizeTask(task, entities.StatusCompleted, "Task completed successfully", startTime)
	case <-task.Ctx.Done():
		wp.tm.logger.Debug().Msgf("Task %s canceled during processing", task.ID)
		return
	case <-wp.tm.stop:
		wp.tm.logger.Debug().Msgf("Service stopped while processing task %s", task.ID)
		return
	}
}

func (wp *workerPool) finalizeTask(task *entities.Task, status entities.TaskStatus, result string, startTime time.Time) {
	now := time.Now()
	task.Status = status
	task.FinishedAt = timePtr(now)
	task.Result = result
	task.Duration = now.Sub(startTime)

	if err := wp.tm.taskRepo.Update(context.Background(), task); err != nil {
		wp.tm.logger.Error().Msgf("failed to update task: %v", err)
	}
}

func (wp *workerPool) wait() {
	wp.wg.Wait()
}

func timePtr(t time.Time) *time.Time {
	return &t
}
