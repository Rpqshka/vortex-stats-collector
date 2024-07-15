package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	vortex "vortex-stats-collector"
)

type getOBInput struct {
	ExchangeName string `json:"exchange_name" binding:"required"`
	Pair         string `json:"pair" binding:"required"`
}

type getOBResponse struct {
	DepthOrder []*vortex.DepthOrder `json:"depth_order"`
}

func (h *Handler) GetOrderBook(c *gin.Context) {
	var input getOBInput

	//Валидация request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	depthOrder, err := h.services.OrderBook.GetOrderBook(input.ExchangeName, input.Pair)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	// Вывод depthOrder из запроса
	c.JSON(http.StatusOK, getOBResponse{
		DepthOrder: depthOrder,
	})
}

type saveOBInput struct {
	ExchangeName string               `json:"exchange_name" binding:"required"`
	Pair         string               `json:"pair" binding:"required"`
	OrderBook    []*vortex.DepthOrder `json:"order_book" binding:"required"`
}

func (h *Handler) SaveOrderBook(c *gin.Context) {
	var input saveOBInput
	//Валидация request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.services.OrderBook.SaveOrderBook(input.ExchangeName, input.Pair, input.OrderBook); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}
