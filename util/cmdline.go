package util

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"

	"github.com/papertrader-api/models"
)

var (
	next      string
	user      models.User
	portfolio models.Portfolio
	game      models.Game
	api       string
	cmdList   = []string{
		"games",
		"orders",
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
	api = fmt.Sprintf(`http://%s:%s`, args["LOCAL"], args["PORT"])
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
	switch s[0] {
	case "games":
		games()
		next = "router"
	case "join":
		join(input)
		next = "router"
	case "buy":
		buy(input)
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
	if err := fetchGame(s[1]); err != nil {
		errorPrompt(fmt.Sprintf("could not fetch game:%s\n", s[1]))
		return
	}
	if err := fetchPortfolio(user.ID, game.ID); err != nil {
		errorPrompt(fmt.Sprintf("could not portfolio:%s\n", s[1]))
		return
	}
	fmt.Printf("game context set\n portfolio id:%s\n", portfolio.ID)
	routerPrompt()
}

func buy(input string) {
	s := strings.Split(input, " ")

	if len(s) < 4 {
		errorPrompt(fmt.Sprintf("unknown input. try:\nsymbol type amount price\n"))
		return
	}
	fmt.Printf("buy order - [0]:%s, [1]:%s, [2]:%s, [3]:%s, [4]:%s\n", s[0], s[1], s[2], s[3], s[4])
	if err := validateContext(); err != nil {
		errorPrompt(err.Error())
		return
	}
	amount, err := strconv.Atoi(s[3])
	if err != nil {
		errorPrompt(fmt.Sprintf("unknown input. count not parse purchase amount\n"))
		return
	}
	if s[2] != "market" && s[4] == "" {
		errorPrompt(fmt.Sprintf("cannot fullfil non market buy order without a price\n"))
		return
	}
	order := models.Order{
		UserID:      user.ID,
		PortfolioID: portfolio.ID,
		OrderAction: models.Buy,
		Symbol:      s[1],
		OrderType:   models.ParseOrderType(s[2]),
		Amount:      amount,
		Ask:         s[4],
	}
	res, err := http.Post(fmt.Sprintf("%s/order/buy", api), "application/json", strings.NewReader(models.ToJson(order)))
	if err != nil {
		errorPrompt(err.Error())
		return
	}
	data, err := ReadResponse(res)
	if err != nil {
		errorPrompt(err.Error())
		return
	}
	json.Unmarshal([]byte(data), &order)
	if order.ID == "" {
		errorPrompt(fmt.Sprintf("error creating order. ID not set"))
		return
	}
	fmt.Printf("order created:\n- id: %s\n- symbol: %s\n- created: %s", order.ID, order.Symbol, order.Created)
}

func errorPrompt(text string) {
	fmt.Printf("%s\n", text)
	routerPrompt()
}

func routerPrompt() {
	fmt.Printf("whats next? %v ", cmdList)
	next = "router"
}

func fetchGame(id string) error {
	apiGame, err := http.Get(fmt.Sprintf("%s/game/%s", api, id))
	if err != nil {
		return err
	}
	gameData, err := ReadResponse(apiGame)
	if err != nil {
		return err
	}
	gameRes := models.APIGameResponse{}
	json.Unmarshal([]byte(gameData), &gameRes)
	if len(gameRes.Data) == 0 {
		return errors.New("coulld not fetch game data")
	}
	game = gameRes.Data[0]
	return nil
}

func fetchPortfolio(userID string, gameID string) error {
	res, err := http.Get(fmt.Sprintf("%s/portfolio/%s", api, gameID))
	if err != nil {
		return err
	}
	data, err := ReadResponse(res)
	if err != nil {
		return err
	}
	portApiRes := models.APIPortfolioResponse{}
	json.Unmarshal([]byte(data), &portApiRes)
	for _, p := range portApiRes.Data {
		if p.UserID == userID {
			portfolio = p
			return nil
		}
	}

	if portfolio.ID == "" {
		if err := generatePortfolio(userID, gameID); err != nil {
			return err
		}
	}

	fmt.Printf("portfolio id: %s\n", portfolio.ID)
	return nil
}

func generatePortfolio(userID string, gameID string) error {
	portfolio = models.Portfolio{
		UserID: userID,
		GameID: gameID,
		Cash:   game.Cap,
		Value:  game.Cap,
		Spent:  "0.00",
	}
	createRes, err := http.Post(fmt.Sprintf("%s/portfolio", api), "application/json", strings.NewReader(models.ToJson(portfolio)))
	if err != nil {
		return err
	}

	data, err := ReadResponse(createRes)
	if err != nil {
		return err
	}
	json.Unmarshal([]byte(data), &portfolio)
	fmt.Printf("portfolio: %s created\n", portfolio.ID)
	return nil
}

func validateContext() error {
	if user.ID == "" {
		return errors.New("user id not set. not sure how you did that")
	}
	if game.ID == "" {
		return errors.New("game id not set. try join game {id} first")
	}
	if portfolio.ID == "" {
		return errors.New("portfolio not set. try join game {id} first")
	}
	return nil
}
