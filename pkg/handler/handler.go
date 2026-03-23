package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grancc/go-effective-mobile-test/pkg/service"
)

type Handler struct {
	services *service.Service
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	return router
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}
