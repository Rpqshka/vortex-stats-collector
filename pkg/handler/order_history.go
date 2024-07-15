package handler

import (
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
	vortex "vortex-stats-collector"
)

type getOHResponse struct {
	HistoryOrder []*vortex.HistoryOrder `json:"history_order"`
}

func (h *Handler) GetOrderHistory(c *gin.Context) {
	var input vortex.Client
	//Валидация request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}

	depthOrder, err := h.services.OrderHistory.GetOrderHistory(&input)
	if err != nil {

		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	// Вывод depthOrder из запроса
	c.JSON(http.StatusOK, getOHResponse{
		HistoryOrder: depthOrder,
	})
}

type saveOHInput struct {
	ClientName          string  `json:"client_name" binding:"required"`
	ExchangeName        string  `json:"exchange_name" binding:"required"`
	Label               string  `json:"label" binding:"required"`
	Pair                string  `json:"pair" binding:"required"`
	Side                string  `json:"side" binding:"required"`
	Type                string  `json:"type" binding:"required"`
	BaseQty             float64 `json:"base_qty" binding:"required"`
	Price               float64 `json:"price" binding:"required"`
	AlgorithmNamePlaced string  `json:"algorithm_name_placed" binding:"required"`
	LowestSellPrc       float64 `json:"lowest_sell_prc" binding:"required"`
	HighestBuyPrc       float64 `json:"highest_buy_prc" binding:"required"`
	CommissionQuoteQty  float64 `json:"commission_quote_qty" binding:"required"`
}

func (h *Handler) SaveOrderHistory(c *gin.Context) {
	var input saveOHInput
	//Валидация request body
	if err := c.BindJSON(&input); err != nil {
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	client := &vortex.Client{
		ClientName:   input.ClientName,
		ExchangeName: input.ExchangeName,
		Label:        input.Label,
		Pair:         input.Pair,
	}
	order := &vortex.HistoryOrder{
		ClientName:          input.ClientName,
		ExchangeName:        input.ExchangeName,
		Label:               input.Label,
		Pair:                input.Pair,
		Side:                input.Side,
		Type:                input.Type,
		BaseQty:             input.BaseQty,
		Price:               input.Price,
		AlgorithmNamePlaced: input.AlgorithmNamePlaced,
		LowestSellPrc:       input.LowestSellPrc,
		HighestBuyPrc:       input.HighestBuyPrc,
		CommissionQuoteQty:  input.CommissionQuoteQty,
		TimePlaced:          time.Now(),
	}
	if err := h.services.OrderHistory.SaveOrder(client, order); err != nil {
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}

	c.JSON(http.StatusOK, map[string]interface{}{
		"status": "success",
	})
}
