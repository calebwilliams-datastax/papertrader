package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/papertrader-api/models"
)

var (
	next      string
	user      models.User
	portfolio models.Portfolio
	api       string
	cmdList   = []string{
		"games",
		"join",
		"buy",
		"sell",
		"stats",
	}
)

func Start(args map[string]string) {
	prompts := map[string]func(string){
		"init":       login,
		"router":     router,
		"usercreate": usercreate,
	}
	next = "init"
	api = fmt.Sprintf(`http://%s:%s`, args["local"], args["port"])
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		os.Exit(0)
	}()
	fmt.Println(`-=papertrader cmd interface=-`)
	fmt.Printf("api url : %s\n", api)
	fmt.Println("whats your username?")
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		input := scanner.Text()
		f := prompts[next]
		f(input)
	}
}
func login(input string) {
	fmt.Printf("fetching user : %s\n", input)
	res, err := http.Get(fmt.Sprintf(`%s/user/%s`, api, input))
	if err != nil {
		errorPrompt(fmt.Sprintf(`error fetching user`))
		return
	}
	usr, err := ReadResponse(res)
	if err != nil {
		errorPrompt(fmt.Sprintf(`error reading user response`))
	}
	apiRes := models.APIUserResponse{}
	json.Unmarshal([]byte(usr), &apiRes)
	if len(apiRes.Data) > 0 {
		user = apiRes.Data[0]
		fmt.Printf("user_id: %s\n", user.ID)
		fmt.Printf("whats next? %v ", cmdList)
		next = "router"
		return
	} else {
		fmt.Printf("type in your username & email ex: foo foo@usernames.com\n(no spaces or stupid characters, please)\n")
		next = "usercreate"
	}
}

func usercreate(input string) {
	s := strings.Split(input, " ")
	//fmt.Printf("creating user - username: %s, email: %s\n", s[0], s[1])
	payload := fmt.Sprintf(`{"name":"%s","email":"%s"}`, s[0], s[1])
	url := fmt.Sprintf(`%s/user`, api)
	fmt.Printf(" post: %s -d %s\n\n", url, payload)
	res, err := http.Post(fmt.Sprintf(`%s/user`, api), `application/json`, strings.NewReader(payload))
	if err != nil {
		errorPrompt(fmt.Sprintf(`could not create user: %s`, err))
	}

	data, err := ReadResponse(res)
	if err != nil {
		errorPrompt(fmt.Sprintf(`error reading response body: %s`, err))
	}
	fmt.Printf("response: %s\n", data)
	fmt.Printf("whats next? %v ", cmdList)
	next = "router"
}

func router(input string) {
	s := strings.Split(input, " ")
	fmt.Printf("router")
	switch s[0] {
	case "games":
		games()
		next = "router"
	case "join":
		join(input)
		next = "router"
	}
}

func games() {
	fmt.Printf("fetching games\n")
	res, err := http.Get(fmt.Sprintf(`%s/game/list`, api))
	if err != nil {
		errorPrompt(fmt.Sprintf(`error fetching game list: %s`, err))
	}
	data, err := ReadResponse(res)
	if err != nil {
		errorPrompt(fmt.Sprintf(`error parsing api response: %s`, err))
	}
	games := models.APIGameResponse{}
	json.Unmarshal([]byte(data), &games)
	fmt.Printf("games count : %v\n========\n", games.Count)
	for _, game := range games.Data {
		fmt.Printf("- id: %s\n- name: %s\n- cap: %v\n- ends: %s\n--------\n", game.ID, game.Name, game.Cap, game.End)
	}
	routerPrompt()
}

func join(input string) {
	fmt.Printf("join - input: %s\n", input)
	s := strings.Split(input, " ")
	if len(s) == 1 {
		errorPrompt("no game id specified. try: join 16172553781ee5\n")
		return
	}
	fmt.Printf("attempting to join game: %s\n", input)
	fmt.Printf("fetching game: %s\n", s[1])
	game, err := fetchGame(s[1])
	if err != nil || game.ID == "" {
		errorPrompt(fmt.Sprintf("could not fetch game:%s\n", s[1]))
		return
	}
	res, err := http.Get(fmt.Sprintf(`%s/portfolio/%s`, api, user.ID))
	if err != nil {
		errorPrompt(fmt.Sprintf(`could not fetch portfolios: %s`, err))
	}
	data, err := ReadResponse(res)
	if err != nil {
		errorPrompt(fmt.Sprintf(`could not parse api response: %s`, err))
	}
	p := models.Portfolio{
		UserID: user.ID,
		GameID: game.ID,
		Cash:   game.Cap,
		Value:  game.Cap,
	}
	portRes := models.APIPortfolioResponse{}
	json.Unmarshal([]byte(data), &portRes)
	for _, portfolio := range portRes.Data {
		if portfolio.GameID == game.ID {
			p = portfolio
		}
	}
	if p.ID == "" {
		portReq, err := http.Post(fmt.Sprintf(`%s/portfolio`, api), `application/json`, strings.NewReader(models.ToJson(p)))
		if err != nil {
			errorPrompt(fmt.Sprintf(`could not create portfolio: %s`, err))
		}
		porRes, err := ReadResponse(portReq)
		if err != nil {
			errorPrompt(fmt.Sprintf(`could not parse portfolio api response: %s`, err))
		}
		fmt.Printf("joined game: %s", porRes)
		routerPrompt()
	}
}

func errorPrompt(text string) {
	fmt.Printf("%s\n", text)
	routerPrompt()
}

func routerPrompt() {
	fmt.Printf("whats next? %v ", cmdList)
	next = "router"
}

func fetchGame(id string) (models.Game, error) {
	game := models.Game{}
	apiGame, err := http.Get(fmt.Sprintf("%s/game/%s", api, s[1]))
	if err != nil {
		return game, err
	}
	gameData, err := ReadResponse(apiGame)
	if err != nil {
		return game, err
	}
	gameRes := models.APIGameResponse{}
	json.Unmarshal([]byte(gameData), &gameRes)
	if len(gameRes.Data) == 0 {
		return game, errors.New("coulld not fetch game data")
	}
	return gameRes.Data[0], nil
}
