package models

import "time"

/*
	games: createTable(
		keyspaceName:"papertrader",
		tableName:"games",
		partitionKeys: [
		  { name: "id", type: {basic: TEXT} }
		],
        values: [
            { name: "cap", type: {basic: DECIMAL }},
            { name: "created", type: {basic: DATE }},
            { name: "created_by", type : {basic: TEXT }},
            { name: "name", type : {basic: TEXT }},
            { name: "end", type: {basic: DATE }}
        ]
	  )
*/

type Game struct {
	ID        string    `json:"id"`
	Created   time.Time `json:"created"`
	CreatedBy string    `json:"created_by"`
	Name      string    `json:"name"`
	End       time.Time `json:"end"`
	Cap       string    `json:"cap"`
}
