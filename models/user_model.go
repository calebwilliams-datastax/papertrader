package models

import "time"

/*
users: createTable(
	  keyspaceName:"papertrader"
	  tableName:"users"
	  partitionKeys: [
		{ name: "id", type: {basic: TEXT} }
        { name: "username", type: {basic: TEXT }}
	  ],
      values: [
          { name: "name", type: {basic: TEXT }}
          { name: "email", type: {basic: TEXT }}
          { name: "created", type: {basic :DATE }}
      ]
	)
*/

type User struct {
	ID      string    `json:"id"`
	Name    string    `json:"name"`
	Email   string    `json:"email"`
	Created time.Time `json:"created"`
}
