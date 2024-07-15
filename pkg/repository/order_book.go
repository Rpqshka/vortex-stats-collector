package repository

import (
	"context"
	"errors"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	vortex "vortex-stats-collector"
)

type OrderBookClickHouse struct {
	db driver.Conn
}

func NewOrderBookClickHouse(db driver.Conn) *OrderBookClickHouse {
	return &OrderBookClickHouse{db: db}
}

// Получение asks и bids
func (r *OrderBookClickHouse) GetOrderBook(exchangeName, pair string) ([]*vortex.DepthOrder, error) {
	query := fmt.Sprintf(`
		SELECT
			arrayMap((x, y) -> (x, y), asks.price, asks.base_qty) as asks,
			arrayMap((x, y) -> (x, y), bids.price, bids.base_qty) as bids
		FROM %s
		WHERE exchange = ? AND pair = ?
	`, orderBookTable)

	rows, err := r.db.Query(context.Background(), query, exchangeName, pair)
	if err != nil {
		return nil, fmt.Errorf("failed to query order book: %w", err)
	}
	defer rows.Close()

	var depthOrders []*vortex.DepthOrder

	for rows.Next() {
		var asks, bids [][]float64

		if err := rows.Scan(&asks, &bids); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		for _, ask := range asks {
			depthOrders = append(depthOrders, &vortex.DepthOrder{
				Price:   ask[0],
				BaseQty: ask[1],
			})
		}

		for _, bid := range bids {
			depthOrders = append(depthOrders, &vortex.DepthOrder{
				Price:   bid[0],
				BaseQty: bid[1],
			})
		}
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during iteration of rows: %w", err)
	}

	return depthOrders, nil
}

func (r *OrderBookClickHouse) SaveOrderBook(exchangeName, pair string, orderBook []*vortex.DepthOrder) error {
	asks, bids, err := splitDepthOrderArray(orderBook)
	if err != nil {
		return err
	}

	askPrices := make([]float64, len(asks))
	askBaseQtys := make([]float64, len(asks))
	bidPrices := make([]float64, len(bids))
	bidBaseQtys := make([]float64, len(bids))

	for i, ask := range asks {
		askPrices[i] = ask.Price
		askBaseQtys[i] = ask.BaseQty
	}

	for i, bid := range bids {
		bidPrices[i] = bid.Price
		bidBaseQtys[i] = bid.BaseQty
	}

	query := fmt.Sprintf(`
			INSERT INTO %s 
			(exchange, pair, asks.price, asks.base_qty, bids.price, bids.base_qty)
			VALUES (?, ?, ?, ?, ?, ?)`, orderBookTable)

	args := []interface{}{
		exchangeName,
		pair,
		askPrices,
		askBaseQtys,
		bidPrices,
		bidBaseQtys,
	}

	if err = r.db.Exec(context.Background(), query, args...); err != nil {
		return err
	}

	return nil
}

// Разделение массива depthOrder на массив asks и bids
func splitDepthOrderArray(depthArr []*vortex.DepthOrder) ([]vortex.DepthOrder, []vortex.DepthOrder, error) {
	if len(depthArr)%2 != 0 {
		return nil, nil, errors.New("the array length is not even")
	}

	mid := len(depthArr) / 2
	asks := make([]vortex.DepthOrder, mid)
	bids := make([]vortex.DepthOrder, mid)
	for i := 0; i < mid; i++ {
		asks[i] = *depthArr[i]
		bids[i] = *depthArr[mid+i]
	}
	return asks, bids, nil

}
