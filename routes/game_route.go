package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (db *EndpointContext) GameById(w http.ResponseWriter, r *http.Request) {

}

func (db *EndpointContext) GameByName(w http.ResponseWriter, r *http.Request) {

}

func (e *EndpointContext) GameCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("GameCreate\n")
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
	w.WriteHeader(code)
	w.Write([]byte(insert))
}

func (db *EndpointContext) GameUpdate(w http.ResponseWriter, r *http.Request) {

}

func (db *EndpointContext) GameDelete(w http.ResponseWriter, r *http.Request) {

}
