package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/handler"
	"vortex-stats-collector/pkg/service"
)

type MockOrderHistoryService struct{}

func (m *MockOrderHistoryService) GetOrderHistory(client *vortex.Client) ([]*vortex.HistoryOrder, error) {
	if client.ClientName == "error" {
		return nil, errors.New("mock GetOrderHistory error")
	}
	return []*vortex.HistoryOrder{
		{
			ClientName:          "client1",
			ExchangeName:        "exchange1",
			Label:               "label1",
			Pair:                "BTC/USD",
			Side:                "buy",
			Type:                "limit",
			BaseQty:             10,
			Price:               100,
			AlgorithmNamePlaced: "algo1",
			LowestSellPrc:       90,
			HighestBuyPrc:       110,
			CommissionQuoteQty:  0.1,
			TimePlaced:          time.Now(),
		},
	}, nil
}

func (m *MockOrderHistoryService) SaveOrder(client *vortex.Client, order *vortex.HistoryOrder) error {
	if client.ClientName == "error" {
		return errors.New("mock SaveOrder error")
	}
	return nil
}

func setupRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/getOrderHistory", handler.GetOrderHistory)
	router.POST("/saveOrderHistory", handler.SaveOrderHistory)

	return router
}

func TestGetOrderHistory(t *testing.T) {
	mockService := &MockOrderHistoryService{}
	services := &service.Service{OrderHistory: mockService}
	handler := handler.NewHandler(services)
	router := setupRouter(handler)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"client_name":   "client1",
				"exchange_name": "exchange1",
				"label":         "label1",
				"pair":          "BTC/USD",
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Invalid request - missing client_name",
			requestBody: map[string]interface{}{
				"exchange_name": "exchange1",
				"label":         "label1",
				"pair":          "BTC/USD",
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request - missing exchange_name",
			requestBody: map[string]interface{}{
				"client_name": "client1",
				"label":       "label1",
				"pair":        "BTC/USD",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("GET", "/getOrderHistory", bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)
			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestSaveOrderHistory(t *testing.T) {
	mockService := &MockOrderHistoryService{}
	services := &service.Service{OrderHistory: mockService}
	handler := handler.NewHandler(services)
	router := setupRouter(handler)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
		expectedError  string
	}{
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"client_name":           "client1",
				"exchange_name":         "exchange1",
				"label":                 "label1",
				"pair":                  "BTC/USD",
				"side":                  "buy",
				"type":                  "limit",
				"base_qty":              10,
				"price":                 100,
				"algorithm_name_placed": "algo1",
				"lowest_sell_prc":       90,
				"highest_buy_prc":       110,
				"commission_quote_qty":  0.1,
			},
			expectedStatus: http.StatusOK,
			expectedError:  "",
		},
		{
			name: "Invalid request - missing client_name",
			requestBody: map[string]interface{}{
				"exchange_name":         "exchange1",
				"label":                 "label1",
				"pair":                  "BTC/USD",
				"side":                  "buy",
				"type":                  "limit",
				"base_qty":              10,
				"price":                 100,
				"algorithm_name_placed": "algo1",
				"lowest_sell_prc":       90,
				"highest_buy_prc":       110,
				"commission_quote_qty":  0.1,
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid request - missing pair",
			requestBody: map[string]interface{}{
				"client_name":           "client2",
				"exchange_name":         "exchange1",
				"label":                 "label1",
				"side":                  "buy",
				"type":                  "limit",
				"base_qty":              10,
				"price":                 100,
				"algorithm_name_placed": "algo1",
				"lowest_sell_prc":       90,
				"highest_buy_prc":       110,
				"commission_quote_qty":  0.1,
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/saveOrderHistory", bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)

		})
	}
}
