package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/calebwilliams-datastax/papertrader-api/models"
	"github.com/calebwilliams-datastax/papertrader-api/util"
	"github.com/gorilla/mux"
)

func (e *EndpointContext) OrdersByPortfolioID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OrdersByPortfolioID\n")
	util.LogRequest(r)
	params := mux.Vars(r)
	code, res, err := e.GetByClause("orders", "$eq", "portfolio_id", []string{params["portfolio_id"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) OrderCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("OrderCreate\n")
	util.LogRequest(r)
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	order := models.Order{}
	json.Unmarshal([]byte(res), &order)
	order.SetDefaults()

	code, dbRes, err := e.PostDB("orders", order)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, dbRes)
	w.WriteHeader(code)

	w.Write([]byte(models.ToJson(order)))
}

func (db *EndpointContext) OrderUpdate(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)
}

func (db *EndpointContext) OrderDelete(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)
}
