package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/papertrader-api/routes"
)

var ec routes.EndpointContext

func main() {
	args := processFlags()
	ec = routes.NewEndpointContext(args)
	log.Printf("papertrader-api starting @ %s:%s", args["local"], args["port"])
	router := mux.NewRouter()
	//healthcheck
	router.HandleFunc("/healthcheck", healthcheck).Methods("GET")
	//setup
	router.HandleFunc("/setup/database/schema", ec.SetupSchema).Methods("GET")
	//users
	router.HandleFunc("/user", ec.UserCreate).Methods("POST")
	router.HandleFunc("/user/{id}", ec.UserById).Methods("GET")
	router.HandleFunc("/user/{name}", ec.UserByName).Methods("GET")
	router.HandleFunc("/user/delete/{id}", ec.UserDelete).Methods("POST")
	//games
	router.HandleFunc("/game", ec.GameCreate).Methods("POST")
	router.HandleFunc("/game/{id}", ec.GameById).Methods("GET")
	router.HandleFunc("/game/{name}", ec.GameByName).Methods("GET")
	router.HandleFunc("/game/delete/{id}", ec.GameDelete).Methods("POST")
	//portfolios
	router.HandleFunc("/portfolio", ec.PortfolioCreate).Methods("POST")
	router.HandleFunc("/portfolio/{id}", ec.PortfolioById).Methods("GET")
	router.HandleFunc("/portfolio/{name}", ec.PortfolioByName).Methods("GET")
	router.HandleFunc("/portfolio/delete/{id}", ec.PortfolioDelete).Methods("POST")
	//orders
	router.HandleFunc("/order/buy", ec.OrderCreate).Methods("POST")
	router.HandleFunc("/order/sell", ec.OrderCreate).Methods("POST")
	router.HandleFunc("/order/{id}", ec.GetOrderById).Methods("GET")
	router.HandleFunc("/order/delete/{id}", ec.PortfolioDelete).Methods("POST")
	//start
	log.Fatal(http.ListenAndServe(fmt.Sprintf("%s:%s", args["local"], args["port"]), router))
}

func processFlags() map[string]string {
	flags := map[string]string{}
	port := flag.String("port", "8080", "-port=8080")
	local := flag.String("local", "localhost", "-local=localhost")
	token := flag.String("token", "637d7289-89a1-4f13-989a-14c6fc837600", "-token=foo")
	db_url := flag.String("db_url", "http://localhost:8080", "-db_url=localhost:8080")
	db_user := flag.String("db_user", "cassandra", "-db_user=cassandra")
	db_pass := flag.String("db_pass", "cassandra", "-db_pass=cassandra")

	flag.Parse()

	flags["port"] = *port
	flags["local"] = *local
	flags["token"] = *token
	flags["db_url"] = *db_url
	flags["db_user"] = *db_user
	flags["db_pass"] = *db_pass
	return flags
}

func healthcheck(w http.ResponseWriter, r *http.Request) {
	log.Printf("healthcheck")
	w.WriteHeader(http.StatusOK)

}
