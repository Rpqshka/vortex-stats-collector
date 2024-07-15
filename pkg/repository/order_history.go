package repository

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
	vortex "vortex-stats-collector"
)

type OrderHistoryClickHouse struct {
	db driver.Conn
}

func NewOrderHistoryClickHouse(db driver.Conn) *OrderHistoryClickHouse {
	return &OrderHistoryClickHouse{db: db}
}

func (r *OrderHistoryClickHouse) GetOrderHistory(client *vortex.Client) ([]*vortex.HistoryOrder, error) {
	query := fmt.Sprintf(`
        SELECT
            client_name,
            exchange_name,
            label,
            pair,
            side,
            type,
            base_qty,
            price,
            algorithm_name_placed,
            lowest_sell_prc,
            highest_buy_prc,
            commission_quote_qty,
            time_placed
        FROM %s
        WHERE client_name = ? AND exchange_name = ? AND label = ? AND pair = ?
    `, orderHistoryTable)

	rows, err := r.db.Query(context.Background(), query,
		client.ClientName, client.ExchangeName, client.Label, client.Pair)
	if err != nil {
		return nil, fmt.Errorf("failed to query order history: %w", err)
	}
	defer rows.Close()

	var orders []*vortex.HistoryOrder

	for rows.Next() {
		var order vortex.HistoryOrder

		if err := rows.Scan(
			&order.ClientName,
			&order.ExchangeName,
			&order.Label,
			&order.Pair,
			&order.Side,
			&order.Type,
			&order.BaseQty,
			&order.Price,
			&order.AlgorithmNamePlaced,
			&order.LowestSellPrc,
			&order.HighestBuyPrc,
			&order.CommissionQuoteQty,
			&order.TimePlaced,
		); err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}

		orders = append(orders, &order)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("failed during iteration of rows: %w", err)
	}

	return orders, nil
}

func (r *OrderHistoryClickHouse) SaveOrder(client *vortex.Client, order *vortex.HistoryOrder) error {
	query := fmt.Sprintf(`
        INSERT INTO order_history (
            client_name, 
            exchange_name, 
            label, 
            pair, 
            side, 
            type, 
            base_qty, 
            price, 
            algorithm_name_placed, 
            lowest_sell_prc, 
            highest_buy_prc, 
            commission_quote_qty, 
            time_placed
        ) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
    `)

	err := r.db.Exec(
		context.Background(),
		query,
		order.ClientName,
		order.ExchangeName,
		order.Label,
		order.Pair,
		order.Side,
		order.Type,
		order.BaseQty,
		order.Price,
		order.AlgorithmNamePlaced,
		order.LowestSellPrc,
		order.HighestBuyPrc,
		order.CommissionQuoteQty,
		order.TimePlaced,
	)
	if err != nil {
		return fmt.Errorf("failed to insert order into ClickHouse: %w", err)
	}

	return nil
}
