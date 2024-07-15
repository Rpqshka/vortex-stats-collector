package service

import (
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/repository"
)

type OrderBook interface {
	GetOrderBook(exchangeName, pair string) ([]*vortex.DepthOrder, error)
	SaveOrderBook(exchangeName, pair string, orderBook []*vortex.DepthOrder) error
}
type OrderHistory interface {
	GetOrderHistory(client *vortex.Client) ([]*vortex.HistoryOrder, error)
	SaveOrder(client *vortex.Client, order *vortex.HistoryOrder) error
}

type Service struct {
	OrderBook
	OrderHistory
}

func NewService(repos *repository.Repository) *Service {
	return &Service{
		OrderBook:    NewOrderBookService(repos.OrderBook),
		OrderHistory: NewOrderHistoryService(repos.OrderHistory),
	}
}
