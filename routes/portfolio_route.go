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

func (e *EndpointContext) PortfolioByUserId(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioByUserId\n")
	params := mux.Vars(r)
	code, res, err := e.GetByClause("portfolios", "$eq", "user_id", []string{params["user_id"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) PortfolioByUserIdGameId(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioByUserIdGameId\n")
	params := mux.Vars(r)
	code, res, err := e.GetByClause("portfolios", "$eq", "user_id", []string{params["user_id"]})
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(res))
}

func (e *EndpointContext) PortfolioCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioCreate\n")
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	portfolio := models.Portfolio{}
	json.Unmarshal([]byte(res), &portfolio)
	portfolio.SetDefaults()

	code, insert, err := e.PostDB("portfolio", portfolio)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(insert))
}

func (e *EndpointContext) PortfolioUpdate(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioUpdate\n")
	res, err := util.ReadRequestBody(r)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	portfolio := models.Portfolio{}
	json.Unmarshal([]byte(res), &portfolio)

	code, dbres, err := e.PutDB("portfolio", portfolio)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(code)
	w.Write([]byte(dbres))
}

func (e *EndpointContext) PortfolioDelete(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioDelete\n")
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
	w.WriteHeader(code)
	w.Write([]byte(dbRes))
}
