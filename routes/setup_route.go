package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/papertrader-api/models"
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

func (ec *EndpointContext) DeleteTestData(w http.ResponseWriter, r *http.Request) {
	ec.RefreshAuthToken()
	query := models.QueryByValues("users", map[string]string{"name": "system"}, models.User{})
	fmt.Printf("query : %s", query)
	res, err := ec.PostGraphQL("keyspace", query)
	if err != nil {
		util.HandleError(w, r, err)
	}
	w.WriteHeader(200)
	w.Write([]byte(res))
}

func (ec *EndpointContext) SetupTestData(w http.ResponseWriter, r *http.Request) {
	usr := models.User{
		ID:      models.GenerateID(),
		Name:    "system",
		Email:   "system@papertrader.com",
		Created: time.Now(),
	}
	game := models.Game{
		ID:        models.GenerateID(),
		Created:   time.Now(),
		CreatedBy: usr.ID,
		End:       time.Now().AddDate(0, 1, 0),
		Name:      "system gen: test game",
		Cap:       "1000000.00",
	}
	portfolio := models.Portfolio{
		ID:      models.GenerateID(),
		Created: time.Now(),
		GameID:  game.ID,
		UserID:  usr.ID,
		Cash:    game.Cap,
		Spent:   "0.00",
	}

	usrQuery := models.InsertMutation("user", "users", usr)
	gameQuery := models.InsertMutation("game", "games", game)
	portfolioQuery := models.InsertMutation("portfolio", "portfolios", portfolio)

	fmt.Printf("===\nuser mut: %s\ngame mut: %s\nport mut:%s\n===", usrQuery, gameQuery, portfolioQuery)
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
mutation kspt { createKeyspace(name:"papertrader", replicas: 1)}`

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
