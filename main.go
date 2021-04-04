package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/routes"
	"github.com/papertrader-api/util"
)

var ec routes.EndpointContext

func main() {
	args := processFlags()
	if args["CMDLINE"] == "true" {
		util.Start(args)
		return
	}
	ec = routes.NewEndpointContext(args)
	log.Printf("papertrader-api starting @ %s:%s", args["LOCAL"], args["PORT"])
	router := mux.NewRouter()
	//healthcheck++
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	//setup
	router.HandleFunc("/setup/database/data", ec.SetupTestData).Methods("GET")
	//users
	router.HandleFunc("/user", ec.UserCreate).Methods("POST")
	router.HandleFunc("/user/{name}", ec.UserByName).Methods("GET")
	router.HandleFunc("/user/delete/{id}", ec.UserDelete).Methods("POST")
	//games
	router.HandleFunc("/game", ec.GameCreate).Methods("POST")
	router.HandleFunc("/game/list", ec.GameList).Methods("GET")
	router.HandleFunc("/game/{id}", ec.GameByID).Methods("GET")
	router.HandleFunc("/game/delete/{id}", ec.GameDelete).Methods("POST")
	//portfolios
	router.HandleFunc("/portfolio", ec.PortfolioCreate).Methods("POST")
	router.HandleFunc("/portfolio/{game_id}", ec.PortfolioByGameID).Methods("GET")
	router.HandleFunc("/portfolio/{game_id}/{user_id}", ec.PortfolioByGameIDUserID).Methods("GET")
	router.HandleFunc("/portfolio/delete/{id}", ec.PortfolioDelete).Methods("POST")
	//orders
	router.HandleFunc("/order/buy", ec.OrderCreate).Methods("POST")
	router.HandleFunc("/order/sell", ec.OrderCreate).Methods("POST")
	router.HandleFunc("/order/{id}", ec.GetOrderByID).Methods("GET")
	router.HandleFunc("/order/delete/{id}", ec.PortfolioDelete).Methods("POST")
	//market
	router.HandleFunc("/market/price/{symbol}", ec.GetPrice).Methods("GET")
	router.HandleFunc("/market/series/{symbol}", ec.TimeSeries).Methods("GET")
	//start
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", args["LOCAL"], args["PORT"]), router))
}

func processFlags() map[string]string {
	flags := map[string]string{}
	port := flag.String("PORT", "8084", "-PORT=8084")
	local := flag.String("LOCAL", "localhost", "-LOCAL=localhost")
	db_url := flag.String("DB_URL", "http://localhost:8082", "-DB_URL=http://localhost:8082")
	auth_url := flag.String("AUTH_URL", "http://localhost:8081", "-AUTH_URL=http://localhost:8081")
	db_user := flag.String("DB_USER", "cassandra", "-DB_USER=cassandra")
	db_pass := flag.String("DB_PASS", "cassandra", "-DB_PASS=cassandra")
	av_url := flag.String("AV_URL", "https://www.alphavantage.co/", "AV_URL=https://www.alphavantage.co/")
	av_token := flag.String("AV_TOKEN", "HR9QB2RM5GWOO0IX", "-AV_TOKEN=foo")
	cmdline := flag.String("CMDLINE", "false", "-CMDLINE=false")

	flag.Parse()

	flags["PORT"] = *port
	flags["LOCAL"] = *local
	flags["DB_URL"] = *db_url
	flags["DB_USER"] = *db_user
	flags["DB_PASS"] = *db_pass
	flags["AUTH_URL"] = *auth_url
	flags["AV_URL"] = *av_url
	flags["AV_TOKEN"] = *av_token
	flags["CMDLINE"] = *cmdline
	return flags
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("healthcheck")
	w.WriteHeader(http.StatusOK)

}
