package routes

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (db *EndpointContext) GetOrderByID(w http.ResponseWriter, r *http.Request) {
	util.LogRequest(r)
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
