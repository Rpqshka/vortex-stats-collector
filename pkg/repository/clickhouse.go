package repository

import (
	"context"
	"fmt"
	"github.com/ClickHouse/clickhouse-go/v2"
	"github.com/ClickHouse/clickhouse-go/v2/lib/driver"
)

// Таблицы БД
const (
	orderBookTable    = "order_book"
	orderHistoryTable = "order_history"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func NewClickhouseDB(cfg Config) (driver.Conn, error) {
	connString := fmt.Sprintf("%s:%s", cfg.Host, cfg.Port)

	options := &clickhouse.Options{
		Addr: []string{connString},
		Auth: clickhouse.Auth{
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.Database,
		},
	}

	conn, err := clickhouse.Open(options)
	if err != nil {
		return nil, fmt.Errorf("failed to open ClickHouse connection: %w", err)
	}

	err = conn.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to ping ClickHouse: %w", err)
	}

	return conn, nil
}
