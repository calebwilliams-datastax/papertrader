package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (ec *EndpointContext) UserById(w http.ResponseWriter, r *http.Request) {

}

func (ec *EndpointContext) UserByName(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("UserByName\n")
	params := mux.Vars(r)
	code, res, err := ec.GetByClause("users", "$eq", "name", []string{params["name"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) UserCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("UserCreate\n")
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	user := models.User{}
	json.Unmarshal([]byte(res), &user)
	user.SetDefaults()

	code, insert, err := e.PostDB("users", user)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(insert))
}

func (ec *EndpointContext) UserUpdate(w http.ResponseWriter, r *http.Request) {}

func (ec *EndpointContext) UserDelete(w http.ResponseWriter, r *http.Request) {}
