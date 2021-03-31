package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

type EndpointContext struct {
	Headers   map[string]string
	Endpoints map[string]string
	Auth      Auth
}

type Auth struct {
	Username string
	Password string
	Token    string
}

func NewEndpointContext(args map[string]string) EndpointContext {
	context := EndpointContext{
		Headers: map[string]string{
			"X-Cassandra-Token": "",
			"Content-Type":      "application/json",
		},
		Endpoints: map[string]string{
			"base":       fmt.Sprintf("%s/", args["db_url"]),
			"auth":       fmt.Sprintf("%s/v1/auth", args["auth_url"]),
			"schema":     fmt.Sprintf("%s/graphql-schema", args["db_url"]),
			"users":      fmt.Sprintf("%s/v2/keyspaces/papertrader/users", args["db_url"]),
			"games":      fmt.Sprintf("%s/v2/keyspaces/papertrader/games", args["db_url"]),
			"portfolios": fmt.Sprintf("%s/v2/keyspaces/papertrader/portfolios", args["db_url"]),
			"orders":     fmt.Sprintf("%s/v2/keyspaces/papertrader/orders", args["db_url"]),
			"keyspace":   fmt.Sprintf("%sv2/schemas/keyspaces", args["db_url"]),
		},
		Auth: Auth{
			Username: args["db_user"],
			Password: args["db_pass"],
		},
	}
	return context
}

func (e *EndpointContext) GetByClause(table, clause, column string, values []string) (int, string, error) {
	e.RefreshAuthToken()
	client := http.Client{}
	where := models.Where(clause, column, values)
	url := fmt.Sprintf(`%s?where=%s`, e.Endpoints[table], where)
	req, err := util.BuildGETRequest(url, e.Headers)
	if err != nil {
		return 500, "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	data, err := util.ReadResponse(res)
	if err != nil {
		return 500, "", err
	}
	return res.StatusCode, data, nil
}

func (e *EndpointContext) PostDB(table string, v interface{}) (int, string, error) {
	e.RefreshAuthToken()
	client := http.Client{}
	defer client.CloseIdleConnections()
	req, err := util.BuildPOSTRequest(e.Endpoints[table],
		models.ToJson(v), e.Headers)
	if err != nil {
		return 500, "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	data, err := util.ReadResponse(res)
	if err != nil {
		return 500, "", err
	}
	return res.StatusCode, data, nil
}

func (e *EndpointContext) PostGraphQL(endpoint string, query string) (string, error) {
	client := http.Client{}
	req, err := http.NewRequest("POST", e.Endpoints[endpoint], strings.NewReader(query))
	if err != nil {
		return "", err
	}
	req.Header.Add("X-Cassandra-Token", e.Auth.Token)
	req.Header.Add("Content-Type", "application/json")
	res, err := client.Do(req)
	if err != nil {
		return "", err
	}
	dbRes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}
	return string(dbRes), nil
}

func (e *EndpointContext) RefreshAuthToken() error {
	res, err := http.Post(e.Endpoints["auth"],
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"username":"%s", "password":"%s"}`, e.Auth.Username, e.Auth.Password))))
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	fmt.Printf("auth res: %s", string(raw))
	data := map[string]string{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	e.Headers["X-Cassandra-Token"] = data["authToken"]
	e.Auth.Token = data["authToken"]
	return nil
}
