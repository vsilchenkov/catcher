package handler

import (
	"catcher/pkg/logging"

	"github.com/gin-gonic/gin"
)

type errorResponse struct {
	Message string `json:"message"`
}

func newErrorResponse(c *gin.Context, statusCode int, op string, err error) {
	logger := logging.GetLogger()
	logger.Error("Error response",
		logger.Str("Request", c.Request.RequestURI),
		logger.Op(op),
		logger.Err(err))
	c.AbortWithStatusJSON(statusCode, errorResponse{err.Error()})
}
