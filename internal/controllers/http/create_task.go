package http

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/controllers/http/middleware"
	"TaskManager/internal/controllers/requests"
	"TaskManager/internal/usecases"
	"TaskManager/pkg/logger"
	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
	"net/http"
)

type createTaskController struct {
	useCase usecases.CreateTaskUseCase
	logger  logger.Logger
}

func NewCreateTaskController(
	handler *gin.Engine,
	useCase usecases.CreateTaskUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	ct := &createTaskController{
		useCase: useCase,
		logger:  logger,
	}

	handler.POST("/task", ct.CreateTask, middleware.HandleErrors)
}

// CreateTask godoc
// @Summary Создание новой задачи
// @Description Создание задачи с помощь передачи текста задания
// @Accept       json
// @Produce      json
// @Param        request body requests.CreateTask true "структура запроса"
// @Success      200 {object} responses.Task
// @Failure      400 {object} string "некорректный формат запроса"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router       /task [post]
func (ct *createTaskController) CreateTask(c *gin.Context) {
	var req requests.CreateTask
	if err := c.ShouldBindJSON(&req); err != nil {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := ct.useCase.CreateTask(c, req)

	if err != nil {
		middleware.AddGinError(c, errors.Wrap(err, "failed to create task"))
		return
	}

	c.JSON(http.StatusOK, response)
}
