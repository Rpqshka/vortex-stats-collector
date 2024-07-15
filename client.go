package vortex_stats_collector

type Client struct {
	ClientName   string `json:"client_name" binding:"required"`
	ExchangeName string `json:"exchange_name" binding:"required"`
	Label        string `json:"label" binding:"required"`
	Pair         string `json:"pair" binding:"required"`
}
