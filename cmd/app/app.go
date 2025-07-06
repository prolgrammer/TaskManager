package app

import (
	"TaskManager/config"
	http2 "TaskManager/internal/controllers/http"
	"TaskManager/internal/controllers/http/middleware"
	"TaskManager/internal/repositories"
	"TaskManager/internal/services"
	"TaskManager/internal/usecases"
	"TaskManager/pkg/logger"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
)

var (
	l logger.Logger

	taskManager services.TaskManager

	createTaskUseCase usecases.CreateTaskUseCase
	getTaskUseCase    usecases.GetTaskUseCase
	deleteTaskUseCase usecases.DeleteTaskUseCase
	getTasksUseCase   usecases.GetTasksUseCase

	taskRepo repositories.TaskRepository
)

func Run() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	err = initPackages(cfg)
	if err != nil {
		log.Fatal(err)
	}

	initRepository(cfg)
	initServices(cfg)
	initUseCases()

	runHTTP(cfg)
}

func initPackages(cfg *config.Config) error {
	var err error

	l = logger.NewConsoleLogger(logger.LevelSwitch(cfg.LogLevel))

	return err
}

func initServices(cfg *config.Config) {
	taskManager = services.NewTaskManager(taskRepo, l, cfg.Services.Workers)
}

func initRepository(cfg *config.Config) {
	taskRepo = repositories.NewTaskRepository()
}

func initUseCases() {
	createTaskUseCase = usecases.NewCreateTaskUseCase(taskRepo, taskManager)
	getTaskUseCase = usecases.NewGetTaskUseCase(taskRepo)
	deleteTaskUseCase = usecases.NewDeleteTaskUseCase(taskRepo)
	getTasksUseCase = usecases.NewGetTasksUseCase(taskRepo)
}

func runHTTP(cfg *config.Config) {
	router := gin.Default()
	router.HandleMethodNotAllowed = true

	mw := middleware.NewMiddleware(l)

	http2.InitMiddleware(router)
	http2.NewCreateTaskController(router, createTaskUseCase, mw, l)
	http2.NewGetTaskController(router, getTaskUseCase, mw, l)
	http2.NewDeleteTaskController(router, deleteTaskUseCase, mw, l)
	http2.NewGetTasksController(router, getTasksUseCase, mw, l)

	address := fmt.Sprintf("%s:%s", cfg.HTTP.Host, cfg.HTTP.Port)
	l.Info().Msgf("starting HTTP server on %s", address)
	err := http.ListenAndServe(address, router)
	if err != nil {
		log.Fatal(err.Error())
	}
}
