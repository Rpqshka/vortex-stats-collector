package repository

import (
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	vortex "vortex-stats-collector"
)

type OrderBook interface {
	GetOrderBook(exchangeName, pair string) ([]*vortex.DepthOrder, error)
	SaveOrderBook(exchangeName, pair string, orderBook []*vortex.DepthOrder) error
}
type OrderHistory interface {
	GetOrderHistory(client *vortex.Client) ([]*vortex.HistoryOrder, error)
	SaveOrder(client *vortex.Client, order *vortex.HistoryOrder) error
}
type Repository struct {
	OrderBook
	OrderHistory
}

func NewRepository(db driver.Conn) *Repository {
	return &Repository{
		OrderBook:    NewOrderBookClickHouse(db),
		OrderHistory: NewOrderHistoryClickHouse(db),
	}
}
