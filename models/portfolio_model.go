package models

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
            {name: "cash", type: {basic: DECIMAL}},
            {name: "spent", type: {basic: DECIMAL }},
            {name: "value", type: { basic: DECIMAL }}
        ]
    )
*/

type Portfolio struct {
	ID     string  `json:"id"`
	GameID string  `json:"game_id"`
	UserID string  `json:"user_id"`
	Cash   float64 `json:"cash"`
	Spent  float64 `json:"spent"`
	Value  float64 `json:"value"`
}
