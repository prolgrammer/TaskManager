package middleware

import (
	"TaskManager/internal/usecases"
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func (m *middleware) HandleErrors(c *gin.Context) {
	if len(c.Errors) > 0 {
		err := c.Errors.Last()

		if errors.Is(err, usecases.ErrEntityAlreadyExists) {
		}

		m.logger.Err(err).Error().Msgf("Unexpected error: ")
		c.AbortWithStatusJSON(http.StatusInternalServerError, "Internal server error")
	}
}
