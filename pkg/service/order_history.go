package service

import (
	vortex "vortex-stats-collector"
	"vortex-stats-collector/pkg/repository"
)

type OrderHistoryService struct {
	repo repository.OrderHistory
}

func NewOrderHistoryService(repo repository.OrderHistory) *OrderHistoryService {
	return &OrderHistoryService{repo: repo}
}

func (s *OrderHistoryService) GetOrderHistory(client *vortex.Client) ([]*vortex.HistoryOrder, error) {
	return s.repo.GetOrderHistory(client)
}

func (s *OrderHistoryService) SaveOrder(client *vortex.Client, order *vortex.HistoryOrder) error {
	return s.repo.SaveOrder(client, order)
}
