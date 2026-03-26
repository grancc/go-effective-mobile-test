package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Message string `json:"message"`
}

func NewErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.WithFields(logrus.Fields{
		"status": statusCode,
		"path":   c.Request.URL.Path,
		"method": c.Request.Method,
	}).Warn(message)
	c.AbortWithStatusJSON(statusCode, Error{message})
}
