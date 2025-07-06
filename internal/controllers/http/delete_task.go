package http

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/controllers/http/middleware"
	"TaskManager/internal/usecases"
	"TaskManager/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
)

type deleteTaskController struct {
	logger  logger.Logger
	useCase usecases.DeleteTaskUseCase
}

func NewDeleteTaskController(
	handler *gin.Engine,
	useCase usecases.DeleteTaskUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	dt := deleteTaskController{
		logger:  logger,
		useCase: useCase,
	}

	handler.DELETE("/task/:task_id", dt.DeleteTask, middleware.HandleErrors)
}

// DeleteTask godoc
// @Summary Удаление задачи
// @Description Удаление задачи по ее id
// @Accept       json
// @Produce      json
// @Param        task_id path string true "path format"
// @Success      200
// @Failure      400 {object} string "некорректный формат запроса"
// @Failure      404 {object} string "задача не найдена"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router       /task/{task_id} [delete]
func (dt *deleteTaskController) DeleteTask(c *gin.Context) {
	taskId := c.Param("task_id")
	if taskId == "" {
		middleware.AddGinError(c, controllers.ErrDataBindError)
		return
	}

	if err := dt.useCase.DeleteTask(c, taskId); err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, "Deleted successful")
}
