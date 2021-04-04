package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/papertrader-api/market"
	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

type EndpointContext struct {
	Headers   map[string]string
	Endpoints map[string]string
	Auth      Auth
	AVConfig  market.AVConfig
}

type Auth struct {
	Username  string
	Password  string
	Token     string
	TokenTime time.Time
}

func NewEndpointContext(args map[string]string) EndpointContext {
	context := EndpointContext{
		Headers: map[string]string{
			"X-Cassandra-Token": "",
			"Content-Type":      "application/json",
		},
		Endpoints: map[string]string{
			"base":       fmt.Sprintf("%s/", args["DB_URL"]),
			"auth":       fmt.Sprintf("%s/v1/auth", args["AUTH_URL"]),
			"schema":     fmt.Sprintf("%s/schema", args["DB_URL"]),
			"users":      fmt.Sprintf("%s/v2/keyspaces/papertrader/users", args["DB_URL"]),
			"games":      fmt.Sprintf("%s/v2/keyspaces/papertrader/games", args["DB_URL"]),
			"portfolios": fmt.Sprintf("%s/v2/keyspaces/papertrader/portfolios", args["DB_URL"]),
			"orders":     fmt.Sprintf("%s/v2/keyspaces/papertrader/orders", args["DB_URL"]),
			"keyspace":   fmt.Sprintf("%sv2/schemas/keyspaces", args["DB_URL"]),
			"delete":     fmt.Sprintf("%s/v2/keyspaces/papertrader", args["DB_URL"]),
		},
		Auth: Auth{
			Username: args["DB_USER"],
			Password: args["DB_PASS"],
		},
		AVConfig: market.AVConfig{
			Base:   args["AV_URL"],
			APIKey: args["AV_TOKEN"],
		},
	}
	if err := context.RefreshAuthToken(); err != nil {
		log.Fatalf(`could not generate authentication token db`)
	}
	return context
}

func (e *EndpointContext) TimeToRefresh() bool {
	expiry := e.Auth.TokenTime.Add(time.Minute * 30)
	now := time.Now()
	if now.Unix() > expiry.Unix() {
		return true
	}
	return false
}

func (e *EndpointContext) GetAll(table string) (int, string, error) {
	if e.TimeToRefresh() {
		e.RefreshAuthToken()
	}
	client := http.Client{}
	url := fmt.Sprintf(`%s?where={}`, e.Endpoints[table])
	req, err := util.BuildGETRequest(url, e.Headers)
	if err != nil {
		return 500, "", err
	}
	res, err := client.Do(req)
	if err != nil {
		return 500, "", err
	}
	data, err := util.ReadResponse(res)
	return res.StatusCode, data, err
}

func (e *EndpointContext) GetByClause(table, clause, column string, values []string) (int, string, error) {
	if e.TimeToRefresh() {
		e.RefreshAuthToken()
	}
	client := http.Client{}
	where := models.Where(clause, column, values)
	url := fmt.Sprintf(`%s?where=%s`, e.Endpoints[table], where)
	fmt.Printf("fetching object: %s, url:%s\n", table, url)
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

func (e *EndpointContext) DeleteDB(table, column, key string) (int, string, error) {
	if e.TimeToRefresh() {
		e.RefreshAuthToken()
	}
	client := http.Client{}
	defer client.CloseIdleConnections()
	url := fmt.Sprintf(`%s/%s/%s`, e.Endpoints["delete"], column, key)
	req, err := util.BuildDELETERequest(url, e.Headers)
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
	if e.TimeToRefresh() {
		e.RefreshAuthToken()
	}
	fmt.Printf("post db : table: %s, obj: %+v", table, v)
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

func (e *EndpointContext) PutDB(table string, v interface{}) (int, string, error) {
	if e.TimeToRefresh() {
		e.RefreshAuthToken()
	}
	client := http.Client{}
	defer client.CloseIdleConnections()
	req, err := util.BuildPUTRequest(e.Endpoints[table],
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
	fmt.Printf("auth res status: %v\n", res.StatusCode)
	data := map[string]string{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	e.Headers["X-Cassandra-Token"] = data["authToken"]
	e.Auth.Token = data["authToken"]
	e.Auth.TokenTime = time.Now()
	return nil
}
