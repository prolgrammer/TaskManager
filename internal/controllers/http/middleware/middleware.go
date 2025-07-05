package middleware

import (
	"TaskManager/pkg/logger"
	"github.com/gin-gonic/gin"
)

type middleware struct {
	logger logger.Logger
}

type Middleware interface {
	HandleErrors(c *gin.Context)
}

func NewMiddleware(
	logger logger.Logger) Middleware {
	return &middleware{logger: logger}
}
