package vortex_stats_collector

import "time"

type OrderBook struct {
	Id       uint         `json:"id"`
	Exchange string       `json:"exchange"`
	Pair     string       `json:"pair"`
	Asks     []DepthOrder `json:"asks"`
	Bids     []DepthOrder `json:"bids"`
}

type DepthOrder struct {
	Price   float64 `json:"price"`
	BaseQty float64 `json:"base_qty"`
}

type HistoryOrder struct {
	ClientName          string    `json:"client_name"`
	ExchangeName        string    `json:"exchange_name"`
	Label               string    `json:"label"`
	Pair                string    `json:"pair"`
	Side                string    `json:"side"`
	Type                string    `json:"type"`
	BaseQty             float64   `json:"base_qty"`
	Price               float64   `json:"price"`
	AlgorithmNamePlaced string    `json:"algorithm_name_placed"`
	LowestSellPrc       float64   `json:"lowest_sell_prc"`
	HighestBuyPrc       float64   `json:"highest_buy_prc"`
	CommissionQuoteQty  float64   `json:"commission_quote_qty"`
	TimePlaced          time.Time `json:"time_placed"`
}
