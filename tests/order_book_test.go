package tests

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/handler"
	"vortex-stats-collector/pkg/service"
)

type MockOrderBookService struct{}

func (m *MockOrderBookService) GetOrderBook(exchangeName, pair string) ([]*vortex.DepthOrder, error) {
	if exchangeName == "error" {
		return nil, errors.New("mock GetOrderBook error")
	}
	return []*vortex.DepthOrder{
		{Price: 100, BaseQty: 10},
		{Price: 110, BaseQty: 20},
	}, nil
}

func (m *MockOrderBookService) SaveOrderBook(exchangeName, pair string, orderBook []*vortex.DepthOrder) error {
	if exchangeName == "error" {
		return errors.New("mock SaveOrderBook error")
	}
	return nil
}

func setupOBRouter(handler *handler.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	router.GET("/getOrderBook", handler.GetOrderBook)
	router.POST("/saveOrderBook", handler.SaveOrderBook)

	return router
}

func TestGetOrderBook(t *testing.T) {
	mockService := &MockOrderBookService{}
	services := &service.Service{OrderBook: mockService}
	handler := handler.NewHandler(services)
	router := setupOBRouter(handler)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"exchange_name": "exchange1",
				"pair":          "BTC/USD",
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request",
			requestBody: map[string]interface{}{
				"exchange_name": "exchange2",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("GET", "/getOrderBook", bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}

func TestSaveOrderBook(t *testing.T) {
	mockService := &MockOrderBookService{}
	services := &service.Service{OrderBook: mockService}
	handler := handler.NewHandler(services)
	router := setupOBRouter(handler)

	testCases := []struct {
		name           string
		requestBody    map[string]interface{}
		expectedStatus int
	}{
		{
			name: "Valid request",
			requestBody: map[string]interface{}{
				"exchange_name": "exchange1",
				"pair":          "BTC/USD",
				"order_book": []map[string]interface{}{
					{"price": 100, "quantity": 10},
					{"price": 110, "quantity": 20},
				},
			},
			expectedStatus: http.StatusOK,
		},
		{
			name: "Invalid request",
			requestBody: map[string]interface{}{
				"exchange_name": "exchange2",
			},
			expectedStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			body, _ := json.Marshal(tc.requestBody)
			req, _ := http.NewRequest("POST", "/saveOrderBook", bytes.NewReader(body))
			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			assert.Equal(t, tc.expectedStatus, w.Code)
		})
	}
}
