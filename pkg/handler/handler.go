package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/grancc/go-effective-mobile-test/pkg/service"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type Handler struct {
	services *service.Service
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.Use(recoveryMiddleware(), logrusMiddleware())

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	api := router.Group("/api")
	{
		subs := api.Group("/subscription")
		{
			subs.POST("/", h.createSubscription)
			subs.GET("/", h.listSubscriptions)
			subs.GET("/subs-sum", h.sumSubscriptions)
			subs.GET("/:id", h.getSubscriptionById)
			subs.PUT("/:id", h.updateSubscriptionById)
			subs.DELETE("/:id", h.deleteSubscriptionById)
		}
	}
	return router
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{services: service}
}
