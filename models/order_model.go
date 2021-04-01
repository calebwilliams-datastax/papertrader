package models

import (
	"time"
)

type APIOrderResponse struct {
	Count int     `json:"count"`
	Data  []Order `json:"data"`
}

type Order struct {
	ID          string      `json:"id"`
	UserID      string      `json:"user_id"`
	PortfolioID string      `json:"portfolio_id"`
	Symbol      string      `json:"symbol"`
	Created     time.Time   `json:"created"`
	Closed      time.Time   `json:"closed"`
	OrderAction OrderAction `json:"action"`
	OrderType   OrderType   `json:"type"`
}

type OrderAction int
type OrderType int

const (
	Buy OrderAction = iota
	Sell
)

const (
	Market OrderType = iota
	Limit
	Stop
	StopLimit
	Trailing
)

func ParseOrderAction(oa string) OrderAction {
	switch oa {
	case "sell":
		return Sell
	default:
		return Buy
	}
}

func ParseOrderType(ot string) OrderType {
	switch ot {
	case "stop":
		return Stop
	case "stoplimit":
		return StopLimit
	case "limit":
		return Limit
	default:
		return Market
	}
}
