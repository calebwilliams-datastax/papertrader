package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (e *EndpointContext) GameList(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("ListGames\n")
	code, res, err := e.GetAll("games")
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) GameByID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GameByID\n")
	util.LogRequest(r)
	params := mux.Vars(r)
	id := params["id"]

	code, res, err := e.GetByClause("games", "$eq", "id", []string{id})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (db *EndpointContext) GameByName(w http.ResponseWriter, r *http.Request) {

}

func (e *EndpointContext) GameCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GameCreate\n")
	util.LogRequest(r)
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	game := models.Game{}
	json.Unmarshal([]byte(res), &game)
	game.SetDefaults()

	code, insert, err := e.PostDB("game", game)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(insert))
}

func (db *EndpointContext) GameUpdate(w http.ResponseWriter, r *http.Request) {

}

func (db *EndpointContext) GameDelete(w http.ResponseWriter, r *http.Request) {

}
