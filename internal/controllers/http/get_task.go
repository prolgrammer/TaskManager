package http

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/controllers/http/middleware"
	"TaskManager/internal/usecases"
	"TaskManager/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type getTaskController struct {
	logger  logger.Logger
	useCase usecases.GetTaskUseCase
}

func NewGetTaskController(
	handler *gin.Engine,
	useCase usecases.GetTaskUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	gt := getTaskController{
		logger:  logger,
		useCase: useCase,
	}

	handler.GET("/task/:task_id", gt.GetTask, middleware.HandleErrors)
}

// GetTask godoc
// @Summary запрос на получение задачи
// @Description запрос на получение задачи с помощью ее ID
// @Produce json
// @Param task_id path string true "id задачи"
// @Success 200 {object} responses.Task
// @Failure 404 {object} string "задача не найдена"
// @Failure 500 {object} string "внутренняя ошибка сервера"
// @Router /task/{task_id} [get]
func (gt *getTaskController) GetTask(c *gin.Context) {
	taskId := c.Param("task_id")
	if taskId == "" {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	response, err := gt.useCase.GetTask(c, taskId)
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
