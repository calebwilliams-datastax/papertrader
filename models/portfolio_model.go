package models

import "time"

/*
portfolios: createTable (
        keyspaceName:"papertrader",
        tableName:"portfolios",
        partitionKeys: [
            {name: "id", type:{basic: TEXT}},
        ],
        clusteringKeys: [
            {name: "game_id", type: {basic: TEXT }},
            {name: "user_id", type: {basic: TEXT }}
        ],
        values: [
            {name:"created", type:{basic: DATE }}
            {name: "cash", type: {basic: DECIMAL}},
            {name: "spent", type: {basic: DECIMAL }},
            {name: "value", type: { basic: DECIMAL }}
        ]
    )
*/

type Portfolio struct {
	ID      string    `json:"id"`
	Created time.Time `json:"created"`
	GameID  string    `json:"game_id"`
	UserID  string    `json:"user_id"`
	Cash    string    `json:"cash"`
	Spent   string    `json:"spent"`
	Value   string    `json:"value"`
}
