package models

import (
	"time"
)

/*
orders: createTable (
        keyspaceName: "papertrader",
        tableName: "orders",
        partitionKeys: [
            {name: "id", type:{basic: TEXT}}
        ],
        clusteringKeys: [
            {name: "user_id", type: {basic: TEXT }},
            {name: "portfolio_id", type: {basic: TEXT }},
        ],
        values: [
            {name: "symbol", type: { basic: TEXT }},
            {name: "action", type: {basic : INT }},
            {name: "created", type: {basic: DATE }},
            {name: "closed", type: {basic: DATE }},
            {name: "type", type: {basic: INT }},
            {name: "ask", type: {basic: DECIMAL }}
            {name: "actual", type:{ basic: DECIMAL }}
        ]
    )
*/

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
