package routes

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/calebwilliams-datastax/papertrader-api/models"
	"github.com/calebwilliams-datastax/papertrader-api/util"
	"github.com/gorilla/mux"
)

func (e *EndpointContext) GameList(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GameList\n")
	code, res, err := e.GetAll("games")
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) GameListOpen(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GameOpenList\n")
	_, res, err := e.GetAll("games")
	if err != nil {
		util.HandleError(w, r, err)
		return
	}

	open := []models.Game{}
	gamesApi := models.APIGameResponse{}
	json.Unmarshal([]byte(res), &gamesApi)
	for _, game := range gamesApi.Data {
		if time.Now().Unix() < game.End.Unix() {
			open = append(open, game)
		}
	}
	gamesApi.Count = len(open)
	gamesApi.Data = open
	data, err := json.Marshal(gamesApi)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(data)
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
