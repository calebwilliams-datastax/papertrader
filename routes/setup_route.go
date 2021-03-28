package routes

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/papertrader-api/util"
)

func (ec *EndpointContext) SetupSchema(w http.ResponseWriter, r *http.Request) {
	if err := ec.RefreshAuthToken(); err != nil {
		util.HandleError(w, r, err)
		return
	}
	response := map[string]int{}
	status, err := ec.CreateKeyspace()
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	response["keyspace"] = status
	status, err = ec.DropTables()
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	response["drop_tables"] = status
	status, err = ec.CreateTables()
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	response["create_tables"] = status
	payload, err := json.Marshal(response)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
	w.Write(payload)
}

func (ec *EndpointContext) CreateKeyspace() (int, error) {
	res, err := http.Post(ec.Endpoints["schema"],
		"application/json",
		bytes.NewBuffer([]byte(createKeyspace)))
	return res.StatusCode, err
}

func (ec *EndpointContext) CreateTables() (int, error) {
	res, err := http.Post(ec.Endpoints["schema"],
		"application/json",
		bytes.NewBuffer([]byte(createTables)))
	return res.StatusCode, err
}

func (ec *EndpointContext) DropTables() (int, error) {
	res, err := http.Post(ec.Endpoints["schema"],
		"application/json",
		bytes.NewBuffer([]byte(dropTables)))
	return res.StatusCode, err
}

var createKeyspace = `
mutation kspt {
    createKeyspace(name:"papertrader", replicas: 1)
}`

var createTables = `
mutation createTables {
	user: createTable(
	  keyspaceName:"papertrader",
	  tableName:"users",
	  partitionKeys: [ 
		{ name: "id", type: {basic: TEXT} }
	  ]
	)
	game: createTable(
		keyspaceName:"papertrader",
		tableName:"games",
		partitionKeys: [ 
		  { name: "id", type: {basic: TEXT} }
		]
	  )
  }`

var dropTables = `
mutation dropTables {
    users: dropTable(
        keyspaceName:"papertrader",
        tableName:"users"
    )
    games: dropTable(
        keyspaceName:"papertrader",
        tableName:"games"
    )
    portfolios: dropTable(
        keyspaceName:"papertrader",
        tableName:"portfolios"
    )
    orders: dropTable(
        keyspaceName:"papertrader",
        tableName:"orders"
    )
}`
