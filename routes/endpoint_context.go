package routes

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type EndpointContext struct {
	Headers   map[string]string
	Endpoints map[string]string
	Auth      Auth
}

type Auth struct {
	Username string
	Password string
}

func NewEndpointContext(args map[string]string) EndpointContext {
	context := EndpointContext{
		Headers: map[string]string{
			"X-Cassandra-Token": args["token"],
		},
		Endpoints: map[string]string{
			"base":       fmt.Sprintf("%s/", args["db_url"]),
			"auth":       fmt.Sprintf("%s/v1/auth", args["db_url"]),
			"schema":     fmt.Sprintf("%s/graphql-schema", args["db_url"]),
			"users":      fmt.Sprintf("%s/graphql/users", args["db_url"]),
			"games":      fmt.Sprintf("%s/graphql/games", args["db_url"]),
			"portfolios": fmt.Sprintf("%s/graphql/portfolios", args["db_url"]),
			"orders":     fmt.Sprintf("%s/graphql/orders", args["db_url"]),
		},
		Auth: Auth{
			Username: args["db_user"],
			Password: args["db_pass"],
		},
	}
	return context
}

func (e *EndpointContext) RefreshAuthToken() error {
	res, err := http.Post(e.Endpoints["auth"],
		"application/json",
		bytes.NewBuffer([]byte(fmt.Sprintf(`{"username":%s, "password":%s}`, e.Auth.Username, e.Auth.Password))))
	if err != nil {
		return err
	}
	raw, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	data := map[string]string{}
	if err := json.Unmarshal(raw, &data); err != nil {
		return err
	}
	e.Headers["X-Cassandra-Token"] = data["authToken"]
	return nil
}
