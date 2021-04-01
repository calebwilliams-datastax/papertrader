package routes

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/market"
	"github.com/papertrader-api/models"
	"github.com/papertrader-api/util"
)

func (e *EndpointContext) GetPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioCreate\n")
	params := mux.Vars(r)
	symbol := params["symbol"]
	if symbol == "" {
		util.HandleError(w, r, errors.New("could not parse symbol"))
		return
	}
	price, err := market.CurrentPrice(symbol, e.AVConfig)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(fmt.Sprintf(`{"symbol":"%s", "price":%v, "time":"%s"}`, symbol, price, time.Now().Local().String())))
}

func (e *EndpointContext) TimeSeries(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("TimeSeries\n")
	params := mux.Vars(r)
	symbol := params["symbol"]
	json.Unmarshal([]byte(symbol), &symbol)
	if symbol == "" {
		util.HandleError(w, r, errors.New("could not parse symbol"))
		return
	}
	series, err := market.TimeSeriesIntraday(symbol, 60, e.AVConfig)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(models.ToJson(series)))
}
