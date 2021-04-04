package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (e *EndpointContext) PortfolioByGameID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioByGameID\n")
	util.LogRequest(r)
	params := mux.Vars(r)
	code, res, err := e.GetByClause("portfolios", "$eq", "game_id", []string{params["game_id"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) PortfolioByGameIDUserID(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioByUserIDGameID\n")
	util.LogRequest(r)
	params := mux.Vars(r)
	code, res, err := e.GetByClause("portfolios", "$eq", "game_id", []string{params["game_id"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) PortfolioCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioCreate\n")
	util.LogRequest(r)
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	portfolio := models.Portfolio{}
	json.Unmarshal([]byte(res), &portfolio)
	portfolio.SetDefaults()

	code, res, err := e.PostDB("portfolios", portfolio)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, res)
	w.WriteHeader(code)
	w.Write([]byte(models.ToJson(portfolio)))
}

func (e *EndpointContext) PortfolioUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioUpdate\n")
	util.LogRequest(r)
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	portfolio := models.Portfolio{}
	json.Unmarshal([]byte(res), &portfolio)

	code, dbRes, err := e.PutDB("portfolio", portfolio)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, dbRes)
	w.WriteHeader(code)
	w.Write([]byte(dbRes))
}

func (e *EndpointContext) PortfolioDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioDelete\n")
	util.LogRequest(r)
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	portfolio := models.Portfolio{}
	json.Unmarshal([]byte(res), &portfolio)
	if portfolio.UserID == "" {
		util.HandleError(w, r, errors.New("could not parse portfolio"))
		return
	}
	code, dbRes, err := e.DeleteDB("portfolios", "user_id", portfolio.UserID)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	util.LogResponse(code, dbRes)
	w.WriteHeader(code)
	w.Write([]byte(dbRes))
}
