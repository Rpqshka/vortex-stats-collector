package service

import (
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/repository"
)

type OrderBookService struct {
	repo repository.OrderBook
}

func NewOrderBookService(repo repository.OrderBook) *OrderBookService {
	return &OrderBookService{repo: repo}
}

func (s *OrderBookService) GetOrderBook(exchangeName, pair string) ([]*vortex.DepthOrder, error) {
	return s.repo.GetOrderBook(exchangeName, pair)
}

func (s *OrderBookService) SaveOrderBook(exchangeName, pair string, orderBook []*vortex.DepthOrder) error {
	return s.repo.SaveOrderBook(exchangeName, pair, orderBook)
}
