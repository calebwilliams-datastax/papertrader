package routes

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/calebwilliams-datastax/papertrader-api/market"
	"github.com/calebwilliams-datastax/papertrader-api/models"
	"github.com/calebwilliams-datastax/papertrader-api/util"
	"github.com/gorilla/mux"
)

func (e *EndpointContext) GetPrice(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("PortfolioCreate\n")
	params := mux.Vars(r)
	symbol := params["symbol"]
	time := params["unix_time"]

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
	w.Write([]byte(fmt.Sprintf(`{"symbol":"%s", "price":%v, "time":"%s"}`, symbol, price, time)))
}

func (e *EndpointContext) TimeSeries(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("TimeSeries\n")
	params := mux.Vars(r)
	symbol := params["symbol"]

	series, err := market.TimeSeriesIntraday(symbol, 60, e.AVConfig)
	if err != nil {
		util.HandleError(w, r, err)
		return
	}
	w.WriteHeader(200)
	w.Write([]byte(models.ToJson(series)))
}
