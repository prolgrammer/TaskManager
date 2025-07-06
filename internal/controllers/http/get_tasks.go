package http

import (
	"TaskManager/internal/controllers"
	"TaskManager/internal/controllers/http/middleware"
	"TaskManager/internal/usecases"
	"TaskManager/pkg/logger"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type getTasksController struct {
	logger  logger.Logger
	useCase usecases.GetTasksUseCase
}

func NewGetTasksController(
	handler *gin.Engine,
	useCase usecases.GetTasksUseCase,
	middleware middleware.Middleware,
	logger logger.Logger,
) {
	gt := getTasksController{
		logger:  logger,
		useCase: useCase,
	}

	handler.GET("/tasks", gt.GetTasks, middleware.HandleErrors)
}

// GetTasks godoc
// @Summary Получение списка задач
// @Description Возвращает список задач с поддержкой пагинации
// @Accept       json
// @Produce      json
// @Param limit query int false "Количество задач на странице" default(10)
// @Param offset query int false "Смещение" default(0)
// @Success      200 {object} []responses.Task
// @Failure      400 {object} string "некорректный формат запроса"
// @Failure      500 {object} string "внутренняя ошибка сервера"
// @Router       /tasks [get]
func (gt *getTasksController) GetTasks(c *gin.Context) {
	limit, err := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if err != nil || limit < 1 {
		middleware.AddGinError(c, controllers.ErrInvalidPaginationParams)
		return
	}

	offset, err := strconv.Atoi(c.DefaultQuery("offset", "0"))
	if err != nil || offset < 0 {
		middleware.AddGinError(c, controllers.ErrInvalidPaginationParams)
		return
	}

	response, err := gt.useCase.GetTasks(c, limit, offset)
	if err != nil {
		middleware.AddGinError(c, err)
		return
	}

	c.JSON(http.StatusOK, response)
}
