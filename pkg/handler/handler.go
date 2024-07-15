package handler

import (
	"github.com/gin-gonic/gin"
	"vortex-stats-collector/pkg/service"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()

	orderBook := router.Group("/order-book")
	{
		orderBook.GET("", h.GetOrderBook)
		orderBook.POST("", h.SaveOrderBook)
	}

	orderHistory := router.Group("/order-history")
	{
		orderHistory.GET("", h.GetOrderHistory)
		orderHistory.POST("", h.SaveOrderHistory)
	}

	return router
}
